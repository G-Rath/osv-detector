package database

import (
	"osv-detector/internal"
)

func (db *ZipDB) Vulnerabilities(includeWithdrawn bool) []OSV {
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

func (db ZipDB) VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) (Vulnerabilities, error) {
	var vulnerabilities Vulnerabilities

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.IsAffected(pkg) && !vulnerabilities.Includes(vulnerability) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities, nil
}

func (db ZipDB) Check(pkgs []internal.PackageDetails) ([]VulnsOrError, error) {
	vulnerabilities := make([]VulnsOrError, 0, len(pkgs))

	for i, pkg := range pkgs {
		vulns, err := db.VulnerabilitiesAffectingPackage(pkg)

		vulnerabilities = append(vulnerabilities, VulnsOrError{
			Index: i,
			Vulns: vulns,
			Err:   err,
		})
	}

	return vulnerabilities, nil
}
