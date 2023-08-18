package database

import (
	"github.com/g-rath/osv-detector/internal"
)

// an OSV database that lives in-memory, and can be used by other structs
// that handle loading the vulnerabilities from where ever
type memDB struct {
	vulnerabilities map[string][]OSV
}

func (db *memDB) addVulnerabilityToMap(osv OSV, hash string) {
	vulns := db.vulnerabilities[hash]

	if vulns == nil {
		vulns = []OSV{}
	}

	db.vulnerabilities[hash] = append(vulns, osv)
}

func (db *memDB) addVulnerability(osv OSV) {
	if len(osv.Affected) == 0 {
		db.addVulnerabilityToMap(osv, "*")
	} else {
		for _, affected := range osv.Affected {
			db.addVulnerabilityToMap(osv, string(affected.Package.Ecosystem)+"-"+affected.Package.NormalizedName())
		}
	}
}

func (db *memDB) Vulnerabilities(includeWithdrawn bool) []OSV {
	var vulnerabilities []OSV

	for _, vulns := range db.vulnerabilities {
		for _, vulnerability := range vulns {
			if (vulnerability.Withdrawn == nil) || includeWithdrawn {
				vulnerabilities = append(vulnerabilities, vulnerability)
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
