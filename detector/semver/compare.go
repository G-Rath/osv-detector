package semver

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

// Compare returns an integer comparing two versions according to
// semantic version precedence.
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
//
// Build versions are ignored.
func (v *Version) Compare(w Version) int {
	return compareComponents(v.Components, w.Components)
}
