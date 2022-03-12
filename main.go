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
	"path"
)

// these come from goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func loadOSVDatabase(offline bool, archiveURL string) database.OSVDatabase {
	db, err := database.NewDB(offline, archiveURL)

	if err != nil {
		msg := err.Error()

		if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
			msg = color.RedString("no local version of the database was found, and --offline flag was set")
		}

		_, _ = fmt.Fprintf(os.Stderr, " %s\n", color.RedString("failed: %s", msg))
		os.Exit(127)
	}

	return *db
}

func printEcosystems(db database.OSVDatabase) {
	ecosystems := db.ListEcosystems()

	fmt.Println("The loaded OSV has vulnerabilities for the following ecosystems:")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s\n", ecosystem)
	}
}

func printPackages(pathToLock string, packages []internal.PackageDetails) {
	fmt.Printf("The following packages were found in %s:\n", pathToLock)

	for _, details := range packages {
		fmt.Printf("  %s: %s@%s\n", details.Ecosystem, details.Name, details.Version)
	}
}

func printVulnerabilities(db database.OSVDatabase, pkg internal.PackageDetails) int {
	vulnerabilities := db.VulnerabilitiesAffectingPackage(pkg)

	if len(vulnerabilities) == 0 {
		return 0
	}

	fmt.Printf(
		"%s %s\n",
		color.YellowString("%s@%s", pkg.Name, pkg.Version),
		color.RedString("is affected by the following vulnerabilities:"),
	)

	for _, vulnerability := range vulnerabilities {
		fmt.Printf(
			"  %s %s\n",
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

func loadEcosystemDatabases(ecosystems []internal.Ecosystem, offline bool) OSVDatabases {
	dbs := make(OSVDatabases, 0, len(ecosystems))

	fmt.Printf("Loading OSV databases for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s", ecosystem)
		archiveURL := ecosystemDatabaseURL(ecosystem)

		db := loadOSVDatabase(offline, archiveURL)

		fmt.Printf(
			" (%s vulnerabilities, including withdrawn - last updated %s)\n",
			color.YellowString("%d", len(db.Vulnerabilities(true))),
			db.UpdatedAt,
		)

		dbs = append(dbs, db)
	}

	fmt.Println()

	return dbs
}

func cacheAllEcosystemDatabases() {
	ecosystems := lockfile.KnownEcosystems()

	loadEcosystemDatabases(ecosystems, false)
}

func main() {
	offline := flag.Bool("offline", false, "Update the OSV database")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to use to determine how to parse the given file")
	printVersion := flag.Bool("version", false, "Print version information")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")
	listPackages := flag.Bool("list-packages", false, "List all the packages that were parsed from the given file")
	cacheAllDatabases := flag.Bool("cache-all-databases", false, "Cache all the known ecosystem databases for offline use")

	flag.Parse()

	if *printVersion {
		fmt.Printf("osv-detector %s (%s, commit %s)\n", version, date, commit)
		os.Exit(0)
	}

	if *cacheAllDatabases {
		cacheAllEcosystemDatabases()
		os.Exit(0)
	}

	pathToLockOrDirectory := flag.Arg(0)

	packages, err := lockfile.Parse(pathToLockOrDirectory, *parseAs)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error, %s\n", err)
		os.Exit(127)
	}

	if *listPackages {
		printPackages(pathToLockOrDirectory, packages)
		os.Exit(0)
	}

	dbs := loadEcosystemDatabases(packages.Ecosystems(), *offline)

	if *listEcosystems {
		for _, db := range dbs {
			printEcosystems(db)
		}
		os.Exit(0)
	}

	file := path.Base(pathToLockOrDirectory)

	knownVulnerabilitiesCount := 0
	for _, pkg := range packages {
		for _, db := range dbs {
			knownVulnerabilitiesCount += printVulnerabilities(db, pkg)
		}
	}

	if knownVulnerabilitiesCount == 0 {
		fmt.Printf("%s\n", color.GreenString("%s has no known vulnerabilities!", file))
		os.Exit(0)
	}

	fmt.Printf("\n%s\n", color.RedString(
		"%s is affected by %d vulnerabilities!",
		file,
		knownVulnerabilitiesCount,
	))

	os.Exit(1)
}
