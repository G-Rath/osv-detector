package semantic

import (
	"github.com/g-rath/osv-detector/internal"
	"math/big"
	"regexp"
)

func convertToBigInt(str string) (*big.Int, bool) {
	i, ok := new(big.Int).SetString(str, 10)

	return i, ok
}

func minInt(x, y int) int {
	if x > y {
		return y
	}

	return x
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}

	return x
}

func compareComponents(a Components, b Components) int {
	numberOfComponents := maxInt(len(a), len(b))

	for i := 0; i < numberOfComponents; i++ {
		diff := a.Fetch(i).Cmp(b.Fetch(i))

		if diff != 0 {
			return diff
		}
	}

	return 0
}

func tryExtractNumber(str string) *big.Int {
	matcher := regexp.MustCompile(`[a-zA-Z.-]+(\d+)`)

	results := matcher.FindStringSubmatch(str)

	if results == nil {
		return big.NewInt(0)
	}

	// it should not be possible for this to not be a number,
	// because we select only numbers above in our regexp
	r, _ := new(big.Int).SetString(results[1], 10)

	return r
}

func compareBuilds(a string, b string) int {
	a = removeBuildMetadata(a)
	b = removeBuildMetadata(b)

	if a == "" && b != "" {
		return +1
	}
	if a != "" && b == "" {
		return -1
	}

	av := tryExtractNumber(a)
	bv := tryExtractNumber(b)

	return av.Cmp(bv)
}

type VersionComparator = func(v, w Version) int

// nolint:gochecknoglobals // this is an optimisation and read-only
var comparators = map[internal.Ecosystem]VersionComparator{
	internal.Ecosystem("npm"):       compareForSemver,
	internal.Ecosystem("crates.io"): compareForSemver,
	internal.Ecosystem("Debian"):    compareForDebian,
	internal.Ecosystem("RubyGems"):  compareForRubyGems,
	internal.Ecosystem("NuGet"):     compareForNuGet,
	internal.Ecosystem("Packagist"): compareForPackagist,
	internal.Ecosystem("Go"):        compareForSemver,
	internal.Ecosystem("Hex"):       compareForSemver,
	internal.Ecosystem("Maven"):     compareForMaven,
	internal.Ecosystem("PyPI"):      compareForPyPI,
}

func compareForFallback(v, w Version) int {
	componentDiff := compareComponents(v.Components, w.Components)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuilds(v.Build, w.Build)
}

func (v *Version) findComparator() VersionComparator {
	if comparator, ok := comparators[v.Ecosystem]; ok {
		return comparator
	}

	return compareForFallback
}

// Compare returns an integer representing the sort order of the given Version w
// relative to the subject Version v.
//
// If the subject Version has an ecosystem, then the comparison will be done in
// accordance to the version specification for that ecosystem (if available);
// otherwise, the comparison will be done using semantic versioning except with
// support for an arbitrary number of components.
//
// In this case, if both versions are considered semantically equal and they both
// have a build string, then a "best effort" comparison will be done, generally by
// attempting to identity a number with the strings and comparing that.
//
// Versions with a build string are considered less than ones without (if both
// have equal components).
//
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
func (v *Version) Compare(w Version) int {
	return v.findComparator()(*v, w)
}

// CompareStr returns an integer representing the sort order of the given Version
// w relative to the subject Version v.
//
// If the subject Version has an ecosystem, then the comparison will be done in
// accordance to the version specification for that ecosystem (if available);
// otherwise, the comparison will be done using semantic versioning except with
// support for an arbitrary number of components.
//
// In this case, if both versions are considered semantically equal and they both
// have a build string, then a "best effort" comparison will be done, generally by
// attempting to identity a number with the strings and comparing that.
//
// Versions with a build string are considered less than ones without (if both
// have equal components).
//
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
func (v *Version) CompareStr(str string) int {
	w := ParseWithEcosystem(str, v.Ecosystem)

	return v.Compare(w)
}
