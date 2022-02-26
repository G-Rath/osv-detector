package main

import (
	"flag"
	"fmt"
	"os"
	"osv-detector/detector/database"
	"osv-detector/detector/parsers"
)

func loadOSVDatabase(offline bool) database.OSVDatabase {
	db, err := database.NewDB(offline, database.GithubOSVDatabaseArchiveURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load the OSV DB: %s\n", err)
		os.Exit(127)
	}

	fmt.Printf("Loaded %d vulnerabilities\n", len(db.Vulnerabilities(true)))

	return *db
}

func printEcosystems(db database.OSVDatabase) {
	ecosystems := db.ListEcosystems()

	fmt.Print("The loaded OSV has vulnerabilities for the following ecosystems:")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s\n", ecosystem)
	}
}

func main() {
	offline := flag.Bool("offline", false, "Update the OSV database")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")

	flag.Parse()
	pathToLockOrDirectory := flag.Arg(0)
	fmt.Println("Hello, world.")

	db := loadOSVDatabase(*offline)

	if *listEcosystems {
		printEcosystems(db)
		os.Exit(0)
	}

	out, err := parsers.ParseComposerLock(pathToLockOrDirectory)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %s\n", pathToLockOrDirectory, err)
		os.Exit(127)
	}

	fmt.Printf("%s", out)
}
