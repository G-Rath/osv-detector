package reporter

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/lockfile"
)

type PackageDetailsWithVulnerabilities struct {
	internal.PackageDetails

	Vulnerabilities database.Vulnerabilities `json:"vulnerabilities"`
	Ignored         database.Vulnerabilities `json:"ignored"`
}

type Report struct {
	lockfile.Lockfile

	// Packages is a map of packages and any vulnerabilities that they're affected by
	Packages []PackageDetailsWithVulnerabilities `json:"packages"`
}

func (r Report) HasKnownVulnerabilities() bool {
	return r.countKnownVulnerabilities() > 0
}

func (r Report) countKnownVulnerabilities() int {
	knownVulnerabilitiesCount := 0

	for _, pkg := range r.Packages {
		knownVulnerabilitiesCount += len(pkg.Vulnerabilities)
	}

	return knownVulnerabilitiesCount
}

func (r Report) HasIgnoredVulnerabilities() bool {
	return r.countIgnoredVulnerabilities() > 0
}

func (r Report) countIgnoredVulnerabilities() int {
	ignoredVulnerabilitiesCount := 0

	for _, pkg := range r.Packages {
		ignoredVulnerabilitiesCount += len(pkg.Ignored)
	}

	return ignoredVulnerabilitiesCount
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

func Form(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}

	return plural
}

func (r Report) describeIgnores() string {
	count := r.countIgnoredVulnerabilities()

	if count == 0 {
		return ""
	}

	return color.YellowString(
		" (%d %s ignored)",
		count,
		Form(count, "was", "were"),
	)
}

func (r Report) String() string {
	count := r.countKnownVulnerabilities()
	ignoreMsg := r.describeIgnores()
	word := "known"

	if ignoreMsg != "" {
		word = "new"
	}

	if count == 0 {
		return fmt.Sprintf(
			"  %s%s\n",
			color.GreenString("no %s vulnerabilities found", word),
			ignoreMsg,
		)
	}

	out := r.formatLineByLine()
	out += "\n"

	out += fmt.Sprintf("\n  %s%s\n",
		color.RedString(
			"%d %s %s found in %s",
			count,
			word,
			Form(count, "vulnerability", "vulnerabilities"),
			r.FilePath,
		),
		ignoreMsg,
	)

	return out
}
