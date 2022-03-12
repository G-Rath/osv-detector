package database

import "osv-detector/internal"

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

type Vulnerabilities []OSV

func (vs Vulnerabilities) Includes(vulnerability OSV) bool {
	for _, osv := range vs {
		if osv.ID == vulnerability.ID {
			return true
		}

		if osv.isAliasOf(vulnerability) {
			return true
		}
		if vulnerability.isAliasOf(osv) {
			return true
		}
	}

	return false
}

func (db *OSVDatabase) VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) Vulnerabilities {
	var vulnerabilities Vulnerabilities

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.IsAffected(pkg) && !vulnerabilities.Includes(vulnerability) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}
