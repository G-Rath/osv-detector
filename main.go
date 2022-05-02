package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"osv-detector/internal"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
	"osv-detector/internal/reporter"
	"path"
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

type OSVDatabases []database.OSVDatabase

func (dbs OSVDatabases) check(lockf lockfile.Lockfile) reporter.Report {
	report := reporter.Report{
		Lockfile: lockf,
		Packages: make([]reporter.PackageDetailsWithVulnerabilities, 0, len(lockf.Packages)),
	}

	for _, pkg := range lockf.Packages {
		vulnerabilities := make(database.Vulnerabilities, 0)

		for _, db := range dbs {
			vulnerabilities = append(vulnerabilities, db.VulnerabilitiesAffectingPackage(pkg)...)
		}

		report.Packages = append(report.Packages, reporter.PackageDetailsWithVulnerabilities{
			PackageDetails:  pkg,
			Vulnerabilities: vulnerabilities,
		})
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

func run() int {
	offline := flag.Bool("offline", false, "Update the OSV database")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to use to determine how to parse the given file")
	printVersion := flag.Bool("version", false, "Print version information")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")
	listPackages := flag.Bool("list-packages", false, "List all the packages that were parsed from the given file")
	cacheAllDatabases := flag.Bool("cache-all-databases", false, "Cache all the known ecosystem databases for offline use")
	outputAsJSON := flag.Bool("json", false, "Output the results in JSON format")

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

	for i, pathToLock := range pathsToLocks {
		if i >= 1 {
			r.PrintText("\n")
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

		dbs, err := loadEcosystemDatabases(r, lockf.Packages.Ecosystems(), *offline)

		if err != nil {
			printDatabaseLoadErr(r, err)
			exitCode = 127

			continue
		}

		report := dbs.check(lockf)

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
