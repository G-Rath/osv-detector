package semver

import (
	"regexp"
	"strconv"
)

func maxInt(x, y int) int {
	if x < y {
		return y
	}

	return x
}

func compareInt(a int, b int) int {
	if a == b {
		return 0
	}

	if a < b {
		return -1
	}

	return +1
}

func compareComponents(a Components, b Components) int {
	numberOfComponents := maxInt(len(a), len(b))

	for i := 0; i < numberOfComponents; i++ {
		diff := compareInt(a.Fetch(i), b.Fetch(i))

		if diff != 0 {
			return diff
		}
	}

	return 0
}

func tryExtractNumber(str string) int {
	matcher := regexp.MustCompile(`[a-zA-Z.-]+(\d+)`)

	results := matcher.FindStringSubmatch(str)

	if results == nil {
		return 0
	}

	r, _ := strconv.Atoi(results[1])

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

	return compareInt(av, bv)
}

// Compare returns an integer comparing two versions according to
// semantic version precedence, then by their build version (if present)
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
//
// Versions with a build are considered less than ones without (if both
// have equal components)
//
// Builds are compared using "best effort" - generally if a build ends with
// a number, that will be used as the main comparator.
func (v *Version) Compare(w Version) int {
	componentDiff := compareComponents(v.Components, w.Components)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuilds(v.Build, w.Build)
}
