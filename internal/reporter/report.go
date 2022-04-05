package reporter

import (
	"fmt"
	"github.com/fatih/color"
	"osv-detector/internal"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
	"strings"
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

func (r Report) formatLineByLine() string {
	lines := make([]string, 0, len(r.Packages))

	for _, pkg := range r.Packages {
		if len(pkg.Vulnerabilities) == 0 {
			continue
		}

		lines = append(lines, fmt.Sprintf(
			"  %s %s",
			color.YellowString("%s@%s", pkg.Name, pkg.Version),
			color.RedString("is affected by the following vulnerabilities:"),
		))

		for _, vulnerability := range pkg.Vulnerabilities {
			lines = append(lines, fmt.Sprintf(
				"    %s %s",
				color.CyanString("%s:", vulnerability.ID),
				vulnerability.Describe(),
			))
		}
	}

	return strings.Join(lines, "\n")
}

func (r Report) ToString() string {
	if r.CountKnownVulnerabilities() == 0 {
		return fmt.Sprintf("%s\n", color.GreenString("  no known vulnerabilities found"))
	}

	out := r.formatLineByLine()
	out += "\n"

	out += fmt.Sprintf("\n  %s\n",
		color.RedString(
			"%d known vulnerabilities found in %s",
			r.CountKnownVulnerabilities(),
			r.FilePath,
		),
	)

	return out
}
