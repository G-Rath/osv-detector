package database

import (
	"encoding/json"
	"fmt"
	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/internal/cachedregexp"
	"github.com/g-rath/osv-detector/pkg/lockfile"
	"github.com/g-rath/osv-detector/pkg/semantic"
	"golang.org/x/exp/slices"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

type AffectsRangeType string

const (
	TypeSemver    AffectsRangeType = "SEMVER"
	TypeEcosystem AffectsRangeType = "ECOSYSTEM"
	TypeGit       AffectsRangeType = "GIT"
)

type Ecosystem = internal.Ecosystem

type Package struct {
	Name      string    `json:"name"`
	Ecosystem Ecosystem `json:"ecosystem"`
}

// NormalizedName ensures that the package name is normalized based on ecosystem
// in accordance to the OSV specification.
//
// This is required because currently both GitHub and Pip seem to be a bit
// inconsistent in their package name handling, so we normalize them
// to be on the safe side.
//
// In the future, it's hoped that this can be improved.
func (p Package) NormalizedName() string {
	if p.Ecosystem != lockfile.PipEcosystem {
		return p.Name
	}

	// per https://www.python.org/dev/peps/pep-0503/#normalized-names
	name := cachedregexp.MustCompile(`[-_.]+`).ReplaceAllString(p.Name, "-")

	return strings.ToLower(name)
}

type RangeEvent struct {
	Introduced   string `json:"introduced,omitempty"`
	Fixed        string `json:"fixed,omitempty"`
	LastAffected string `json:"last_affected,omitempty"`
}

type AffectsRange struct {
	Type   AffectsRangeType `json:"type"`
	Events []RangeEvent     `json:"events"`
}

func (e RangeEvent) version() string {
	if e.Introduced != "" {
		return e.Introduced
	}

	if e.Fixed != "" {
		return e.Fixed
	}

	if e.LastAffected != "" {
		return e.LastAffected
	}

	return ""
}

func (ar AffectsRange) containsVersion(pkg internal.PackageDetails) bool {
	if ar.Type != TypeEcosystem && ar.Type != TypeSemver {
		return false
	}
	// todo: we should probably warn here
	if len(ar.Events) == 0 {
		return false
	}

	vp := semantic.MustParse(pkg.Version, pkg.CompareAs)

	sort.Slice(ar.Events, func(i, j int) bool {
		a := ar.Events[i]
		b := ar.Events[j]

		if a.Introduced == "0" {
			return true
		}

		if b.Introduced == "0" {
			return false
		}

		return semantic.MustParse(a.version(), pkg.CompareAs).CompareStr(b.version()) < 0
	})

	var affected bool
	for _, e := range ar.Events {
		if affected {
			if e.Fixed != "" {
				affected = vp.CompareStr(e.Fixed) < 0
			} else if e.LastAffected != "" {
				affected = e.LastAffected == pkg.Version || vp.CompareStr(e.LastAffected) <= 0
			}
		} else if e.Introduced != "" {
			affected = e.Introduced == "0" || vp.CompareStr(e.Introduced) >= 0
		}
	}

	return affected
}

type Affects []AffectsRange

// affectsVersion checks if the given version is within the range
// specified by the events of any "Ecosystem" or "Semver" type ranges
func (a Affects) affectsVersion(pkg internal.PackageDetails) bool {
	for _, r := range a {
		if r.Type != TypeEcosystem && r.Type != TypeSemver {
			return false
		}
		if r.containsVersion(pkg) {
			return true
		}
	}

	return false
}

type Reference struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Versions []string

// MarshalJSON ensures that if there are no versions,
// an empty array is used as the value instead of "null"
func (vs Versions) MarshalJSON() ([]byte, error) {
	if len(vs) == 0 {
		return []byte("[]"), nil
	}

	out, err := json.Marshal([]string(vs))

	if err != nil {
		return out, fmt.Errorf("%w", err)
	}

	return out, nil
}

type Affected struct {
	Package  Package  `json:"package"`
	Versions Versions `json:"versions"`
	Ranges   Affects  `json:"ranges,omitempty"`
}

// OSV represents an OSV style JSON vulnerability database entry
type OSV struct {
	ID        string     `json:"id"`
	Aliases   []string   `json:"aliases"`
	Summary   string     `json:"summary"`
	Published time.Time  `json:"published"`
	Modified  time.Time  `json:"modified"`
	Withdrawn *time.Time `json:"withdrawn,omitempty"`
	Details   string     `json:"details"`
	Affected  []Affected `json:"affected"`
}

func (osv *OSV) isAliasOf(vulnerability OSV) bool {
	for _, alias := range vulnerability.Aliases {
		if osv.ID == alias || slices.Contains(osv.Aliases, alias) {
			return true
		}
	}

	return false
}

func (osv *OSV) AffectsEcosystem(ecosystem internal.Ecosystem) bool {
	for _, affected := range osv.Affected {
		if affected.Package.Ecosystem == ecosystem {
			return true
		}
	}

	return false
}

// truncate ensures that the given string is shorter than the provided limit.
//
// If the string is longer than the limit, it's trimmed and suffixed with an ellipsis.
// Ideally the string will be trimmed at the space that's closest to the limit to
// preserve whole words; if a string has no spaces before the limit, it'll be forcefully truncated.
func truncate(str string, limit int) string {
	count := 0
	truncateAt := -1

	for i, c := range str {
		if unicode.IsSpace(c) {
			truncateAt = i
		}

		count++

		if count >= limit {
			// ideally we want to keep words whole when truncating,
			// but if we can't find a space just truncate at the limit
			if truncateAt == -1 {
				truncateAt = limit
			}

			return str[:truncateAt] + "..."
		}
	}

	return str
}

func (osv *OSV) Describe() string {
	description := osv.Summary

	if description == "" {
		description += truncate(osv.Details, 80)
	}

	if description == "" {
		description += "(no details available)"
	}

	if link := osv.Link(); link != "" {
		description += " (" + link + ")"
	}

	return description
}

// Link returns a URL to the advisory, if possible.
// Otherwise, an empty string is returned
func (osv *OSV) Link() string {
	if strings.HasPrefix(osv.ID, "GHSA") {
		return "https://github.com/advisories/" + osv.ID
	}

	return ""
}

func (osv *OSV) IsAffected(pkg internal.PackageDetails) bool {
	for _, affected := range osv.Affected {
		if affected.Package.Ecosystem == pkg.Ecosystem &&
			affected.Package.NormalizedName() == pkg.Name {
			if len(affected.Ranges) == 0 && len(affected.Versions) == 0 {
				_, _ = fmt.Fprintf(
					os.Stderr,
					"%s does not have any ranges or versions - this is probably a mistake!\n",
					osv.ID,
				)

				continue
			}

			if slices.Contains(affected.Versions, pkg.Version) {
				return true
			}

			if affected.Ranges.affectsVersion(pkg) {
				return true
			}

			if pkg.Version == "" {
				return true
			}
		}
	}

	return false
}
