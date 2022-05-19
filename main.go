package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"osv-detector/internal"
	"osv-detector/internal/configer"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
	"osv-detector/internal/reporter"
	"path"
	"sort"
)

// these come from goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printDatabaseLoadErr(r *reporter.Reporter, err error) {
	msg := err.Error()

	if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		msg = color.RedString("no local version of the database was found, and --offline flag was set")
	}

	r.PrintError(fmt.Sprintf(" %s\n", color.RedString("failed: %s", msg)))
}

func printKnownEcosystems(r *reporter.Reporter) {
	ecosystems := lockfile.KnownEcosystems()

	r.PrintText("The detector supports parsing for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		r.PrintText(fmt.Sprintf("  %s\n", ecosystem))
	}
}

func ecosystemDatabaseURL(ecosystem internal.Ecosystem) string {
	return fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", ecosystem)
}

type VulnsOrError = struct {
	database.Vulnerabilities
	error
}

type GenericOSVDatabase2 interface {
	Check(pkgs []internal.PackageDetails) []database.Vulnerabilities
	// Check(pkgs []internal.PackageDetails) ([]VulnsOrError, error)
}

type GenericOSVDatabase interface {
	Check(pkgs []internal.PackageDetails) []database.Vulnerabilities
	VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) database.Vulnerabilities
}

type GenericOSVDatabaseWithBulk interface {
	VulnerabilitiesAffectingPackages(pkgs lockfile.Packages) []database.Vulnerabilities
}

type OSVDatabases []GenericOSVDatabase

func contains(items []string, value string) bool {
	for _, item := range items {
		if value == item {
			return true
		}
	}

	return false
}

// a struct to hold the result from each request including an index
// which will be used for sorting the results after they come in
type result struct {
	index int
	res   reporter.PackageDetailsWithVulnerabilities
	err   error
}

func (dbs OSVDatabases) parallelFetch2(pkgs lockfile.Packages, ignores []string, concurrencyLimit int) []result {
	if concurrencyLimit > 1 {
		return dbs.parallelFetch(pkgs, ignores, concurrencyLimit)
	}

	var results []result

	if len(pkgs) == 0 {
		return results
	}

	vulns := make([][]database.Vulnerabilities, 0, len(dbs))

	for _, db := range dbs {
		vs := db.Check(pkgs)

		if len(vs) != len(pkgs) {
			panic("oh noes!")
		}

		vulns = append(vulns, vs)
	}

	for i, pkg := range pkgs {
		results = append(results, result{
			index: i,
			res:   dbs.transposePkgResults(pkg, ignores, i, vulns),
			err:   nil,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	return results
}

func (dbs OSVDatabases) parallelFetch(pkgs lockfile.Packages, ignores []string, concurrencyLimit int) []result {
	var results []result

	if len(pkgs) == 0 {
		return results
	}

	// buffered channel which controls the number of concurrent operations
	semaphoreChan := make(chan struct{}, concurrencyLimit)
	resultsChan := make(chan *result)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, pkg := range pkgs {
		go func(i int, pkg lockfile.PackageDetails) {
			// read from the buffered semaphore channel, which will block if we're
			// already got as many goroutines as our concurrency limit allows
			//
			// when one of those routines finish they'll read from this channel,
			// freeing up a slot to unblock this send
			semaphoreChan <- struct{}{}

			vulns := dbs.checkPkg(pkg, ignores)
			result := &result{i, vulns, nil}

			resultsChan <- result

			// read from the buffered semaphore to free up a slot to allow
			// another goroutine to start, since this one is wrapping up
			<-semaphoreChan
		}(i, pkg)
	}

	for {
		result := <-resultsChan
		results = append(results, *result)

		if len(results) == len(pkgs) {
			break
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	return results
}

func (dbs OSVDatabases) transposePkgResults(
	pkg internal.PackageDetails,
	ignores []string,
	packageIndex int,
	allVulns [][]database.Vulnerabilities,
) reporter.PackageDetailsWithVulnerabilities {
	vulnerabilities := make(database.Vulnerabilities, 0)
	ignored := make(database.Vulnerabilities, 0)

	for _, vulns1 := range allVulns {
		vulns := vulns1[packageIndex]

		for _, vulnerability := range vulns {
			if vulnerabilities.Includes(vulnerability) || ignored.Includes(vulnerability) {
				continue
			}

			if contains(ignores, vulnerability.ID) {
				ignored = append(ignored, vulnerability)
			} else {
				vulnerabilities = append(vulnerabilities, vulnerability)
			}
		}
	}

	return reporter.PackageDetailsWithVulnerabilities{
		PackageDetails:  pkg,
		Vulnerabilities: vulnerabilities,
		Ignored:         ignored,
	}
}

func (dbs OSVDatabases) checkPkg(pkg internal.PackageDetails, ignores []string) reporter.PackageDetailsWithVulnerabilities {
	vulnerabilities := make(database.Vulnerabilities, 0)
	ignored := make(database.Vulnerabilities, 0)

	for i, db := range dbs {
		for _, vulnerability := range db.VulnerabilitiesAffectingPackage(pkg) {
			if vulnerabilities.Includes(vulnerability) || ignored.Includes(vulnerability) {
				continue
			}

			if i > 0 {
				fmt.Printf("oh noes!")
			}

			if contains(ignores, vulnerability.ID) {
				ignored = append(ignored, vulnerability)
			} else {
				vulnerabilities = append(vulnerabilities, vulnerability)
			}
		}
	}

	return reporter.PackageDetailsWithVulnerabilities{
		PackageDetails:  pkg,
		Vulnerabilities: vulnerabilities,
		Ignored:         ignored,
	}
}

func (dbs OSVDatabases) check(lockf lockfile.Lockfile, ignores []string, conLimit int) reporter.Report {
	report := reporter.Report{
		Lockfile: lockf,
		Packages: make([]reporter.PackageDetailsWithVulnerabilities, 0, len(lockf.Packages)),
	}

	for _, r := range dbs.parallelFetch2(lockf.Packages, ignores, conLimit) {
		if r.err != nil {
			fmt.Printf("error happened: %v\n", r.err)

			continue
		}

		report.Packages = append(report.Packages, r.res)
	}

	return report
}

func loadEcosystemDatabases(r *reporter.Reporter, ecosystems []internal.Ecosystem, offline bool) (OSVDatabases, error) {
	dbs := make(OSVDatabases, 0, len(ecosystems))

	r.PrintText("  Loading OSV databases for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		r.PrintText(fmt.Sprintf("    %s", ecosystem))
		archiveURL := ecosystemDatabaseURL(ecosystem)

		db, err := database.NewDB(offline, archiveURL)

		if err != nil {
			return dbs, fmt.Errorf("could not load database: %w", err)
		}

		count := len(db.Vulnerabilities(true))

		r.PrintText(fmt.Sprintf(
			" (%s %s, including withdrawn - last updated %s)\n",
			color.YellowString("%d", count),
			reporter.Form(count, "vulnerability", "vulnerabilities"),
			db.UpdatedAt,
		))

		dbs = append(dbs, *db)
	}

	r.PrintText("\n")

	return dbs, nil
}

func cacheAllEcosystemDatabases(r *reporter.Reporter) error {
	ecosystems := lockfile.KnownEcosystems()

	_, err := loadEcosystemDatabases(r, ecosystems, false)

	return err
}

func findLockfiles(r *reporter.Reporter, pathToLockOrDirectory string, parseAs string) []string {
	lockfiles := make([]string, 0, 1)
	file, err := os.Open(pathToLockOrDirectory)

	if err == nil {
		info, err := file.Stat()

		if err == nil {
			if info.IsDir() {
				dirs, err := file.ReadDir(-1)

				if err == nil {
					for _, dir := range dirs {
						if dir.IsDir() {
							continue
						}

						if p, _ := lockfile.FindParser(dir.Name(), parseAs); p == nil {
							continue
						}

						lockfiles = append(lockfiles, path.Join(pathToLockOrDirectory, dir.Name()))
					}
				}
			} else {
				lockfiles = append(lockfiles, pathToLockOrDirectory)
			}
		}
	}

	if err != nil {
		r.PrintError(fmt.Sprintf("Error reading %s: %v\n", pathToLockOrDirectory, err))
	}

	return lockfiles
}

func findAllLockfiles(r *reporter.Reporter, pathsToCheck []string, parseAs string) []string {
	var paths []string

	for _, pathToLockOrDirectory := range pathsToCheck {
		paths = append(paths, findLockfiles(r, pathToLockOrDirectory, parseAs)...)
	}

	return paths
}

type stringsFlag []string

func (s *stringsFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)

	return nil
}

func allIgnores(global, local []string) []string {
	ignores := make(
		[]string,
		0,
		// len cannot return negative numbers, but the types can't reflect that
		uint64(len(global))+uint64(len(local)),
	)

	ignores = append(ignores, global...)
	ignores = append(ignores, local...)

	return ignores
}

func run() int {
	var ignores stringsFlag

	offline := flag.Bool("offline", false, "Perform checks using only the cached databases on disk")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to parse the input files as")
	configPath := flag.String("config", "", "Path to a config file to use for all lockfiles")
	noConfig := flag.Bool("no-config", false, "Disable loading of any config files")
	printVersion := flag.Bool("version", false, "Print version information")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all of the known ecosystems that are supported by the detector")
	listPackages := flag.Bool("list-packages", false, "List the packages that are parsed from the input files")
	cacheAllDatabases := flag.Bool("cache-all-databases", false, "Cache all the known ecosystem databases for offline use")
	outputAsJSON := flag.Bool("json", false, "Output the results in JSON format")
	useDatabases := flag.Bool("use-dbs", false, "Output the results in JSON format")
	useAPI := flag.Bool("use-api", false, "Output the results in JSON format")
	conLimit := flag.Int("n", 1, "Output the results in JSON format")

	flag.Var(&ignores, "ignore", `ID of an OSV to ignore when determining exit codes.
This flag can be passed multiple times to ignore different vulnerabilities`)

	flag.Parse()

	r := reporter.New(os.Stdout, os.Stderr, *outputAsJSON)
	if *outputAsJSON {
		defer r.PrintJSONResults()
	}

	if *printVersion {
		r.PrintText(fmt.Sprintf("osv-detector %s (%s, commit %s)\n", version, date, commit))

		return 0
	}

	if *cacheAllDatabases {
		err := cacheAllEcosystemDatabases(r)

		if err != nil {
			printDatabaseLoadErr(r, err)

			return 127
		}

		return 0
	}

	if *listEcosystems {
		printKnownEcosystems(r)

		return 0
	}

	if *parseAs != "" {
		if parser, parsedAs := lockfile.FindParser("", *parseAs); parser == nil {
			r.PrintError(fmt.Sprintf("Don't know how to parse files as \"%s\" - supported values are:\n", parsedAs))

			for _, s := range lockfile.ListParsers() {
				r.PrintError(fmt.Sprintf("  %s\n", s))
			}

			return 127
		}
	}

	pathsToLocks := findAllLockfiles(r, flag.Args(), *parseAs)

	if len(pathsToLocks) == 0 {
		r.PrintError(
			"You must provide at least one path to either a lockfile or a directory containing a lockfile (see --help for usage and flags)\n",
		)

		return 127
	}

	exitCode := 0

	var config configer.Config

	if !*noConfig && *configPath != "" {
		con, err := configer.Load(*configPath)

		if err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", err))

			return 127
		}

		config = con
	}

	for i, pathToLock := range pathsToLocks {
		config := config

		if i >= 1 {
			r.PrintText("\n")
		}

		if !*noConfig && *configPath == "" {
			base := path.Dir(pathToLock)
			con, err := configer.Find(base)

			if err != nil {
				r.PrintError(fmt.Sprintf("Error, %s\n", err))
				exitCode = 127

				continue
			}

			config = con
		}

		lockf, err := lockfile.Parse(pathToLock, *parseAs)

		if err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", err))
			exitCode = 127

			continue
		}

		r.PrintText(fmt.Sprintf(
			"%s: found %s %s\n",
			color.MagentaString("%s", lockf.FilePath),
			color.YellowString("%d", len(lockf.Packages)),
			reporter.Form(len(lockf.Packages), "package", "packages"),
		))

		if *listPackages {
			r.PrintResult(lockf)

			continue
		}

		// an empty FilePath means we didn't load a config
		if config.FilePath != "" {
			r.PrintText(fmt.Sprintf(
				"  Using config at %s (%s)\n",
				color.MagentaString(config.FilePath),
				color.YellowString("%d %s",
					len(config.Ignore),
					reporter.Form(len(config.Ignore), "ignore", "ignores"),
				),
			))
		}

		var dbs OSVDatabases

		if *useDatabases {
			loaded, err := loadEcosystemDatabases(r, lockf.Packages.Ecosystems(), *offline)

			if err != nil {
				printDatabaseLoadErr(r, err)
				exitCode = 127

				continue
			}

			dbs = append(dbs, loaded...)
		}

		if *useAPI {
			dbs = append(dbs, database.APIOSVDatabase{})
		}

		report := dbs.check(lockf, allIgnores(config.Ignore, ignores), *conLimit)

		r.PrintResult(report)

		if report.HasKnownVulnerabilities() && exitCode == 0 {
			exitCode = 1
		}
	}

	return exitCode
}

func main() {
	os.Exit(run())
}
