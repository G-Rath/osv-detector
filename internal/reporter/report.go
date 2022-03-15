package reporter

import (
	"osv-detector/internal"
	"osv-detector/internal/database"
)

type PackageDetailsWithVulnerabilities struct {
	Name      string
	Version   string
	Ecosystem internal.Ecosystem

	Vulnerabilities database.Vulnerabilities
}

type Report struct {
	FilePath string
	ParsedAs string
	// Packages is a map of packages and any vulnerabilities that they're affected by
	Packages []PackageDetailsWithVulnerabilities
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
