package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"osv-detector/detector"
	"osv-detector/detector/database"
	"osv-detector/detector/parsers"
	"path"
)

func loadOSVDatabase(offline bool) database.OSVDatabase {
	db, err := database.NewDB(offline, database.GithubOSVDatabaseArchiveURL)

	if err != nil {
		msg := fmt.Sprintf("Error loading the OSV DB: %s", err)

		if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
			msg = "Error: --offline can only be used when a local version of the OSV database is available"
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n", msg)
		os.Exit(127)
	}

	fmt.Printf(
		"Loaded %s vulnerabilities (including withdrawn, last updated %s)\n",
		color.YellowString("%d", len(db.Vulnerabilities(true))),
		db.UpdatedAt,
	)

	return *db
}

func printEcosystems(db database.OSVDatabase) {
	ecosystems := db.ListEcosystems()

	fmt.Print("The loaded OSV has vulnerabilities for the following ecosystems:")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s\n", ecosystem)
	}
}

func printPackages(pathToLock string, packages []detector.PackageDetails) {
	fmt.Printf("The following packages were found in %s:\n", pathToLock)

	for _, details := range packages {
		fmt.Printf("  %s: %s@%s\n", details.Ecosystem, details.Name, details.Version)
	}
}

func printVulnerabilities(db database.OSVDatabase, pkg detector.PackageDetails) int {
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
			vulnerability.Summary,
		)
	}

	return len(vulnerabilities)
}

func main() {
	offline := flag.Bool("offline", false, "Update the OSV database")
	parseAs := flag.String("parse-as", "", "Name of a supported lockfile to use to determine how to parse the given file")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")
	listPackages := flag.Bool("list-packages", false, "List all the packages that were parsed from the given file")

	flag.Parse()
	pathToLockOrDirectory := flag.Arg(0)

	db := loadOSVDatabase(*offline)

	if *listEcosystems {
		printEcosystems(db)
		os.Exit(0)
	}

	packages, err := parsers.TryParse(pathToLockOrDirectory, *parseAs)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error, %s\n", err)
		os.Exit(127)
	}

	if *listPackages {
		printPackages(pathToLockOrDirectory, packages)
		os.Exit(0)
	}

	file := path.Base(pathToLockOrDirectory)

	knownVulnerabilitiesCount := 0
	for _, pkg := range packages {
		knownVulnerabilitiesCount += printVulnerabilities(db, pkg)
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
