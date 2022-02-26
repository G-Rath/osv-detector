package database

import (
	"osv-detector/detector"
)

func (db *OSVDatabase) Vulnerabilities(includeWithdrawn bool) []OSV {
	if includeWithdrawn {
		return db.vulnerabilities
	}

	var vulnerabilities []OSV

	for _, vulnerability := range db.vulnerabilities {
		if vulnerability.Withdrawn == nil {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}

func (db *OSVDatabase) VulnerabilitiesAffectingPackage(pkg detector.PackageDetails) []OSV {
	var vulnerabilities []OSV

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.IsAffected(pkg) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}
