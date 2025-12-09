package database

import (
	"github.com/g-rath/osv-detector/internal"
)

// an OSV database that lives in-memory, and can be used by other structs
// that handle loading the vulnerabilities from where ever
type memDB struct {
	vulnerabilities      map[string][]OSV
	VulnerabilitiesCount int
}

func (db *memDB) addVulnerability(osv OSV, pkgNames []string) {
	db.VulnerabilitiesCount++

	// if we have been provided a list of package names, only load advisories
	// that might actually affect those packages, rather than all advisories
	if len(pkgNames) != 0 && !mightAffectPackages(osv, pkgNames) {
		return
	}

	for _, affected := range osv.Affected {
		hash := string(affected.Package.NormalizedEcosystem()) + "-" + affected.Package.NormalizedName()
		vulns := db.vulnerabilities[hash]

		if vulns == nil {
			vulns = []OSV{}
		}

		db.vulnerabilities[hash] = append(vulns, osv)
	}
}

func (db *memDB) Vulnerabilities(includeWithdrawn bool) []OSV {
	var vulnerabilities []OSV
	ids := make(map[string]struct{})

	for _, vulns := range db.vulnerabilities {
		for _, vulnerability := range vulns {
			if _, ok := ids[vulnerability.ID]; ok {
				continue
			}

			if (vulnerability.Withdrawn == nil) || includeWithdrawn {
				vulnerabilities = append(vulnerabilities, vulnerability)
				ids[vulnerability.ID] = struct{}{}
			}
		}
	}

	return vulnerabilities
}

func (db *memDB) VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) Vulnerabilities {
	var vulnerabilities Vulnerabilities

	hash := string(pkg.Ecosystem) + "-" + pkg.Name

	if vulns, ok := db.vulnerabilities[hash]; ok {
		for _, vulnerability := range vulns {
			if vulnerability.Withdrawn == nil && vulnerability.IsAffected(pkg) && !vulnerabilities.Includes(vulnerability) {
				vulnerabilities = append(vulnerabilities, vulnerability)
			}
		}
	}

	return vulnerabilities
}

func (db *memDB) Check(pkgs []internal.PackageDetails) ([]Vulnerabilities, error) {
	vulnerabilities := make([]Vulnerabilities, 0, len(pkgs))

	for _, pkg := range pkgs {
		vulnerabilities = append(vulnerabilities, db.VulnerabilitiesAffectingPackage(pkg))
	}

	return vulnerabilities, nil
}
