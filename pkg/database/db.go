package database

import "osv-detector/internal"

type DB interface {
	Name() string

	// Identifier can be used to check what config this database represents
	Identifier() string

	// Check looks for known vulnerabilities for the given pkgs within this OSV database.
	//
	// The vulnerabilities are returned in an array whose index align with the index of
	// the package that they're for within the pkgs array that was given.
	Check(pkgs []internal.PackageDetails) ([]Vulnerabilities, error)
}
