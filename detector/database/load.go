package database

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// GithubOSVDatabaseArchiveURL represents GitHub's OSV database URL
const GithubOSVDatabaseArchiveURL = "https://codeload.github.com/github/advisory-database/zip/main"

// load fetches a zip archive of the OSV database and loads known vulnerabilities
// from it (which are assumed to be in json files following the OSV spec).
//
// Internally, the archive is cached along with the date that it was fetched
// so that a new version of the archive is only downloaded if it has been
// modified, per HTTP caching standards.
func (db *OSVDatabase) load() error {
	db.vulnerabilities = []OSV{}

	cache, err := db.fetchCache()

	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(cache.Body), int64(len(cache.Body)))
	if err != nil {
		return err
	}

	// Read all the files from the zip archive
	for _, zipFile := range zipReader.File {
		// todo: somehow support passing these excludes generically, so it's more agnostic
		if !strings.HasPrefix(zipFile.Name, "advisory-database-main/advisories") {
			continue
		}

		if strings.HasPrefix(zipFile.Name, "advisory-database-main/advisories/unreviewed") {
			continue
		}

		if !strings.HasSuffix(zipFile.Name, ".json") {
			continue
		}

		file, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not read %s: %v", zipFile.Name, err)

			continue
		}

		var pa OSV
		if err := json.Unmarshal(content, &pa); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s is not a valid JSON file: %v", zipFile.Name, err)

			continue
		}

		db.vulnerabilities = append(db.vulnerabilities, pa)
	}

	return nil
}
