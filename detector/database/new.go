package database

import "fmt"

// OSVDatabase stores security advisories in the manner defined by the OSV spec
type OSVDatabase struct {
	vulnerabilities []OSV
	ArchiveURL      string
	Offline         bool
	UpdatedAt       string
}

// NewDB creates an OSV database with vulnerabilities loaded from an archive
func NewDB(offline bool, dbArchiveURL string) (*OSVDatabase, error) {
	db := &OSVDatabase{Offline: offline, ArchiveURL: dbArchiveURL}
	if err := db.load(); err != nil {
		return nil, fmt.Errorf("unable to fetch OSV database: %w", err)
	}

	return db, nil
}
