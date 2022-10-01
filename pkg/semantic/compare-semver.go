package semantic

import (
	"fmt"
	"strings"
)

func compareNumericComponents(a Components, b Components) int {
	numberOfComponents := maxInt(len(a), len(b))

	for i := 0; i < numberOfComponents; i++ {
		diff := a.Fetch(i).Cmp(b.Fetch(i))

		if diff != 0 {
			return diff
		}
	}

	return 0
}

// Removes build metadata from the given string if present, per semver v2
//
// See https://semver.org/spec/v2.0.0.html#spec-item-10
func removeBuildMetadata(str string) string {
	parts := strings.Split(str, "+")

	return parts[0]
}

func compareBuildComponents(a, b string) int {
	// https://semver.org/spec/v2.0.0.html#spec-item-10
	a = removeBuildMetadata(a)
	b = removeBuildMetadata(b)

	// the spec doesn't explicitly say "don't include the hyphen in the compare"
	// but it's what node-semver does so for now let's go with that...
	a = strings.TrimPrefix(a, "-")
	b = strings.TrimPrefix(b, "-")

	// versions with a prerelease are considered less than those without
	// https://semver.org/spec/v2.0.0.html#spec-item-9
	if a == "" && b != "" {
		return +1
	}
	if a != "" && b == "" {
		return -1
	}

	return compareSemverBuildComponents(
		strings.Split(a, "."),
		strings.Split(b, "."),
	)
}

func compareSemverBuildComponents(a, b []string) int {
	min := minInt(len(a), len(b))

	var compare int

	for i := 0; i < min; i++ {
		ai, aIsNumber := convertToBigInt(a[i])
		bi, bIsNumber := convertToBigInt(b[i])

		switch {
		// 1. Identifiers consisting of only digits are compared numerically.
		case aIsNumber && bIsNumber:
			compare = ai.Cmp(bi)
		// 2. Identifiers with letters or hyphens are compared lexically in ASCII sort order.
		case !aIsNumber && !bIsNumber:
			compare = strings.Compare(a[i], b[i])
		// 3. Numeric identifiers always have lower precedence than non-numeric identifiers.
		case aIsNumber:
			compare = -1
		default:
			compare = +1
		}

		if compare != 0 {
			if compare > 0 {
				return 1
			}

			return -1
		}
	}

	// 4. A larger set of pre-release fields has a higher precedence than a smaller set,
	//    if all the preceding identifiers are equal.
	if len(a) > len(b) {
		return +1
	}
	if len(a) < len(b) {
		return -1
	}

	return 0
}

func (v *SemverLikeVersion) fetchComponentsAndBuild(maxComponents int) (Components, string) {
	if len(v.Components) <= maxComponents {
		return v.Components, v.Build
	}

	comps := v.Components[:maxComponents]
	extra := v.Components[maxComponents:]

	build := v.Build

	for _, c := range extra {
		build += fmt.Sprintf(".%d", c)
	}

	return comps, build
}

type SemverVersion struct {
	SemverLikeVersion
}

func parseSemverVersion(str string) SemverVersion {
	return SemverVersion{ParseSemverLikeVersion(str, 3)}
}

func (v SemverVersion) Compare(w SemverVersion) int {
	componentDiff := compareNumericComponents(v.Components, w.Components)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuildComponents(v.Build, w.Build)
}

func (v SemverVersion) CompareStr(str string) int {
	return v.Compare(parseSemverVersion(str))
}
