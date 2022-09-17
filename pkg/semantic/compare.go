package semantic

import (
	"github.com/g-rath/osv-detector/pkg/lockfile"
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

// Compare returns an integer comparing two versions according to the rules of
// the left-hand versions ecosystem if set, otherwise falling back to semantic
// version precedence and then by their build version (if present).
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
//
// Versions with a build are considered less than ones without (if both
// have equal components)
//
// Builds are compared using "best effort" - generally if a build ends with
// a number, that will be used as the main comparator.
func (v *Version) Compare(w Version) int {
	if v.Ecosystem == lockfile.ComposerEcosystem {
		return compareForPackagist(v.OriginStr, w.OriginStr)
	}

	componentDiff := compareComponents(v.Components, w.Components)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuilds(v.Build, w.Build)
}

// CompareStr returns an integer comparing two versions according to the rules of
// the left-hand versions ecosystem if set, otherwise falling back to semantic
// version precedence and then by their build version (if present).
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
//
// Versions with a build are considered less than ones without (if both
// have equal components)
//
// Builds are compared using "best effort" - generally if a build ends with
// a number, that will be used as the main comparator.
func (v *Version) CompareStr(str string) int {
	w := ParseWithEcosystem(str, v.Ecosystem)

	return v.Compare(w)
}
