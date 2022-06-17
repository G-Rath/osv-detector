package database

import (
	"osv-detector/internal"
)

func (db *DirDB) Vulnerabilities(includeWithdrawn bool) []OSV {
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

func (db DirDB) VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) Vulnerabilities {
	var vulnerabilities Vulnerabilities

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.IsAffected(pkg) && !vulnerabilities.Includes(vulnerability) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}

func (db DirDB) Check(pkgs []internal.PackageDetails) ([]Vulnerabilities, error) {
	vulnerabilities := make([]Vulnerabilities, 0, len(pkgs))

	for _, pkg := range pkgs {
		vulnerabilities = append(vulnerabilities, db.VulnerabilitiesAffectingPackage(pkg))
	}

	return vulnerabilities, nil
}
