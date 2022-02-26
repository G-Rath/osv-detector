package database

import "fmt"

// OSVDatabase stores security advisories in the manner defined by the OSV spec
type OSVDatabase struct {
	vulnerabilities []OSV
	ArchiveURL      string
	Offline         bool
}

// todo: support settings, including "way to load database"
// e.g. we want to set rules for how to exclude json files
// that way it'll be more agnostic

// NewDB fetches the advisory DB from GitHub
func NewDB(offline bool, dbArchiveURL string) (*OSVDatabase, error) {
	db := &OSVDatabase{Offline: offline, ArchiveURL: dbArchiveURL}
	if err := db.load(); err != nil {
		return nil, fmt.Errorf("unable to fetch OSV database: %w", err)
	}

	return db, nil
}
