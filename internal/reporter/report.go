package reporter

import (
	"osv-detector/internal"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
)

type PackageDetailsWithVulnerabilities struct {
	internal.PackageDetails

	Vulnerabilities database.Vulnerabilities `json:"vulnerabilities"`
}

type Report struct {
	lockfile.Lockfile
	// Packages is a map of packages and any vulnerabilities that they're affected by
	Packages []PackageDetailsWithVulnerabilities `json:"packages"`
}

func (r Report) CountKnownVulnerabilities() int {
	knownVulnerabilitiesCount := 0

	for _, pkg := range r.Packages {
		knownVulnerabilitiesCount += len(pkg.Vulnerabilities)
	}

	return knownVulnerabilitiesCount
}

func (r Report) Format(asJSON bool) string {
	if asJSON {
		return r.FormatJSON()
	}

	return r.FormatLineByLine()
}
