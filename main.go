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
)

// these come from goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printDatabaseLoadErr(err error) int {
	msg := err.Error()

	if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		msg = color.RedString("no local version of the database was found, and --offline flag was set")
	}

	_, _ = fmt.Fprintf(os.Stderr, " %s\n", color.RedString("failed: %s", msg))

	return 127
}

func printKnownEcosystems() {
	ecosystems := lockfile.KnownEcosystems()

	fmt.Println("The detector supports parsing for the following ecosystems:")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s\n", ecosystem)
	}
}

func printPackages(lockf lockfile.Lockfile) {
	fmt.Printf("The following packages were found in %s:\n", lockf.FilePath)

	for _, details := range lockf.Packages {
		fmt.Printf("  %s: %s@%s\n", details.Ecosystem, details.Name, details.Version)
	}
}

func printVulnerabilities(db database.OSVDatabase, pkg internal.PackageDetails) int {
	vulnerabilities := db.VulnerabilitiesAffectingPackage(pkg)

	if len(vulnerabilities) == 0 {
		return 0
	}

	fmt.Printf(
		"  %s %s\n",
		color.YellowString("%s@%s", pkg.Name, pkg.Version),
		color.RedString("is affected by the following vulnerabilities:"),
	)

	for _, vulnerability := range vulnerabilities {
		fmt.Printf(
			"    %s %s\n",
			color.CyanString("%s:", vulnerability.ID),
			vulnerability.Describe(),
		)
	}

	return len(vulnerabilities)
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

func loadEcosystemDatabases(ecosystems []internal.Ecosystem, offline bool) (OSVDatabases, error) {
	dbs := make(OSVDatabases, 0, len(ecosystems))

	fmt.Fprintf(os.Stderr, "  Loading OSV databases for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		fmt.Fprintf(os.Stderr, "    %s", ecosystem)
		archiveURL := ecosystemDatabaseURL(ecosystem)

		db, err := database.NewDB(offline, archiveURL)

		if err != nil {
			return dbs, fmt.Errorf("could not load database: %w", err)
		}

		fmt.Fprintf(os.Stderr,
			" (%s vulnerabilities, including withdrawn - last updated %s)\n",
			color.YellowString("%d", len(db.Vulnerabilities(true))),
			db.UpdatedAt,
		)

		dbs = append(dbs, *db)
	}

	fmt.Fprintln(os.Stderr)

	return dbs, nil
}

func cacheAllEcosystemDatabases() error {
	ecosystems := lockfile.KnownEcosystems()

	_, err := loadEcosystemDatabases(ecosystems, false)

	return err
}

func run() int {
	offline := flag.Bool("offline", false, "Update the OSV database")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to use to determine how to parse the given file")
	printVersion := flag.Bool("version", false, "Print version information")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")
	listPackages := flag.Bool("list-packages", false, "List all the packages that were parsed from the given file")
	cacheAllDatabases := flag.Bool("cache-all-databases", false, "Cache all the known ecosystem databases for offline use")
	outputAsJSON := flag.Bool("json", false, "Cache all the known ecosystem databases for offline use")

	flag.Parse()

	if *printVersion {
		fmt.Printf("osv-detector %s (%s, commit %s)\n", version, date, commit)

		return 0
	}

	if *cacheAllDatabases {
		err := cacheAllEcosystemDatabases()

		if err != nil {
			return printDatabaseLoadErr(err)
		}

		return 0
	}

	if *listEcosystems {
		printKnownEcosystems()

		return 0
	}

	pathsToCheck := flag.Args()

	if len(pathsToCheck) == 0 {
		fmt.Fprintf(os.Stderr, "Error, no package information found (see --help for usage)")

		return 1
	}

	exitCode := 0

	for i, pathToLockOrDirectory := range pathsToCheck {
		if i >= 1 {
			fmt.Println()
		}

		lockf, err := lockfile.Parse(pathToLockOrDirectory, *parseAs)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error, %s\n", err)
			exitCode = 127

			continue
		}

		fmt.Fprintf(os.Stderr,
			"%s: found %s packages\n",
			color.MagentaString("%s", lockf.FilePath),
			color.YellowString("%d", len(lockf.Packages)),
		)

		if *listPackages {
			printPackages(lockf)
			continue
		}

		dbs, err := loadEcosystemDatabases(lockf.Packages.Ecosystems(), *offline)

		if err != nil {
			exitCode = printDatabaseLoadErr(err)

			continue
		}

		report := dbs.check(lockf)

		if report.CountKnownVulnerabilities() == 0 {
			fmt.Printf("%s\n", color.GreenString("  no known vulnerabilities found"))

			continue
		}

		fmt.Println(report.Format(*outputAsJSON))

		fmt.Fprintf(os.Stderr,
			"\n  %s\n",
			color.RedString(
				"%d known vulnerabilities found in %s",
				report.CountKnownVulnerabilities(),
				pathToLockOrDirectory,
			),
		)

		if exitCode == 0 {
			exitCode = 1
		}
	}

	return exitCode
}

func main() {
	os.Exit(run())
}
