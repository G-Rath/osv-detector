package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"osv-detector/internal"
	"osv-detector/internal/configer"
	"osv-detector/internal/reporter"
	"osv-detector/pkg/database"
	"osv-detector/pkg/lockfile"
	"path"
	"strings"
)

// these come from goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func ecosystemDatabaseURL(ecosystem internal.Ecosystem) string {
	return fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", ecosystem)
}

type OSVDatabases []database.DB

func contains(items []string, value string) bool {
	for _, item := range items {
		if value == item {
			return true
		}
	}

	return false
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
			// skip vulnerabilities that were already included from a previous database
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

func (dbs OSVDatabases) check(r *reporter.Reporter, lockf lockfile.Lockfile, ignores []string) reporter.Report {
	report := reporter.Report{
		Lockfile: lockf,
		Packages: make([]reporter.PackageDetailsWithVulnerabilities, 0, len(lockf.Packages)),
	}

	vulns := make([][]database.Vulnerabilities, 0, len(dbs))

	for _, db := range dbs {
		results, err := db.Check(lockf.Packages)

		if err != nil {
			r.PrintError(color.RedString(fmt.Sprintf(
				"  an api error occurred while trying to check the packages listed in %s: %v\n",
				lockf.FilePath,
				err,
			)))

			continue
		}

		vulns = append(vulns, results)
	}

	for i, pkg := range lockf.Packages {
		report.Packages = append(
			report.Packages,
			dbs.transposePkgResults(pkg, ignores, i, vulns),
		)
	}

	return report
}

func loadEcosystemDatabases(r *reporter.Reporter, ecosystems []internal.Ecosystem, offline bool) (OSVDatabases, error) {
	dbs := make(OSVDatabases, 0, len(ecosystems))

	r.PrintText("Loading OSV databases for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		r.PrintText(fmt.Sprintf("  %s", ecosystem))
		archiveURL := ecosystemDatabaseURL(ecosystem)

		db, err := database.NewZippedDB(archiveURL, offline)

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

func handleParseAsCSV(r *reporter.Reporter, lines []string, offline bool, ignores []string) int {
	if len(lines) == 0 {
		r.PrintError("You must provide at least one CSV line to parse\n")

		return 127
	}

	lockf, err := lockfile.FromCSV(strings.Join(lines, "\n"))

	if err != nil {
		r.PrintError(fmt.Sprintf("Error, %s\n", err))

		return 127
	}

	r.PrintText(fmt.Sprintf(
		"%s: found %s packages\n",
		color.MagentaString("%s", lockf.FilePath),
		color.YellowString("%d", len(lockf.Packages)),
	))

	dbs, err := loadEcosystemDatabases(r, lockf.Packages.Ecosystems(), offline)

	if err != nil {
		printDatabaseLoadErr(r, err)

		return 127
	}

	report := dbs.check(lockf, ignores)

	r.PrintResult(report)

	if report.HasKnownVulnerabilities() {
		return 1
	}

	return 0
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

type lockfileAndConfigOrErr struct {
	lockf  lockfile.Lockfile
	config configer.Config
	err    error
}

func readAllLockfiles(pathsToLocks []string, parseAs string, checkForLocalConfig bool, config configer.Config) []lockfileAndConfigOrErr {
	lockfiles := make([]lockfileAndConfigOrErr, 0, len(pathsToLocks))

	for _, pathToLock := range pathsToLocks {
		if checkForLocalConfig {
			base := path.Dir(pathToLock)
			con, err := configer.Find(base)

			if err != nil {
				lockfiles = append(lockfiles, lockfileAndConfigOrErr{lockfile.Lockfile{}, config, err})

				continue
			}

			config = con
		}

		lockf, err := lockfile.Parse(pathToLock, parseAs)
		lockfiles = append(lockfiles, lockfileAndConfigOrErr{lockf, config, err})
	}

	return lockfiles
}

func collectEcosystems(files []lockfileAndConfigOrErr) []internal.Ecosystem {
	var ecosystems []internal.Ecosystem

	for _, result := range files {
		if result.err != nil {
			continue
		}

		for _, ecosystem := range result.lockf.Packages.Ecosystems() {
			alreadyExists := false

			for _, eco := range ecosystems {
				if alreadyExists {
					continue
				}

				if eco == ecosystem {
					alreadyExists = true
				}
			}

			if alreadyExists {
				continue
			}

			ecosystems = append(ecosystems, ecosystem)
		}
	}

	return ecosystems
}

func loadDatabases(
	r *reporter.Reporter,
	ecosystems []internal.Ecosystem,
	useDatabases bool,
	useAPI bool,
	batchSize int,
	offline bool,
) (OSVDatabases, bool) {
	var dbs OSVDatabases
	errored := false

	if useDatabases {
		loaded, err := loadEcosystemDatabases(r, ecosystems, offline)

		if err != nil {
			r.PrintDatabaseLoadErr(err)
			errored = true
		} else {
			dbs = append(dbs, loaded...)
		}
	}

	if useAPI {
		db, err := database.NewAPIDB("https://api.osv.dev/v1", batchSize, offline)

		if err != nil {
			r.PrintDatabaseLoadErr(err)
			errored = true
		} else {
			dbs = append(dbs, db)
		}
	}

	return dbs, errored
}

func run() int {
	var ignores stringsFlag

	offline := flag.Bool("offline", false, "Perform checks using only the cached databases on disk")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to parse the input files as")
	parseAsCSV := flag.Bool("parse-as-csv", false, "Parse the input as CSV rows and files")
	configPath := flag.String("config", "", "Path to a config file to use for all lockfiles")
	noConfig := flag.Bool("no-config", false, "Disable loading of any config files")
	printVersion := flag.Bool("version", false, "Print version information")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all of the known ecosystems that are supported by the detector")
	listPackages := flag.Bool("list-packages", false, "List the packages that are parsed from the input files")
	cacheAllDatabases := flag.Bool("cache-all-databases", false, "Cache all the known ecosystem databases for offline use")
	outputAsJSON := flag.Bool("json", false, "Output the results in JSON format")
	useDatabases := flag.Bool("use-dbs", true, "Use the databases from osv.dev to check for known vulnerabilities")
	useAPI := flag.Bool("use-api", false, "Use the osv.dev API to check for known vulnerabilities")
	batchSize := flag.Int("batch-size", 1000, "The number of packages to include in each batch when using the api database")

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
			r.PrintDatabaseLoadErr(err)

			return 127
		}

		return 0
	}

	if *listEcosystems {
		r.PrintKnownEcosystems()

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

	if *parseAsCSV {
		return handleParseAsCSV(r, flag.Args(), *offline, ignores)
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
	loadLocalConfig := !*noConfig

	if loadLocalConfig && *configPath != "" {
		con, err := configer.Load(*configPath)

		if err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", err))

			return 127
		}

		config = con
		loadLocalConfig = false
	}

	files := readAllLockfiles(pathsToLocks, *parseAs, loadLocalConfig, config)

	ecosystems := collectEcosystems(files)

	dbs, errored := loadDatabases(
		r,
		ecosystems,
		*useDatabases,
		*useAPI,
		*batchSize,
		*offline,
	)

	if errored {
		exitCode = 127
	}

	for i, result := range files {
		if i >= 1 {
			r.PrintText("\n")
		}

		if result.err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", result.err))
			exitCode = 127

			continue
		}

		config := result.config
		lockf := result.lockf

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

		report := dbs.check(r, lockf, allIgnores(config.Ignore, ignores))

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
