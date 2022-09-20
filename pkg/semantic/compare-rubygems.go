package semantic

import (
	"strconv"
	"strings"
)

func canonicalizeRubyGemVersion(str string) string {
	res := ""

	checkPrevious := false
	previousWasDigit := true

	for _, c := range str {
		if c == 46 {
			checkPrevious = false
			res += "."

			continue
		}

		isDigit := c >= 48 && c <= 57

		if checkPrevious && previousWasDigit != isDigit {
			res += "."
		}

		res += string(c)

		previousWasDigit = isDigit
		checkPrevious = true
	}

	return res
}

func compareRubyGemsComponents(a, b []string) int {
	min := minInt(len(a), len(b))

	var compare int

	for i := 0; i < min; i++ {
		ai, aIsNumber := convertToBigInt(a[i])
		bi, bIsNumber := convertToBigInt(b[i])

		switch {
		case aIsNumber && bIsNumber:
			compare = ai.Cmp(bi)
		case !aIsNumber && !bIsNumber:
			compare = strings.Compare(a[i], b[i])
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

	if len(a) > len(b) {
		next := a[len(b)]

		if _, err := strconv.Atoi(next); err == nil {
			return 1
		}

		return -1
	}

	if len(a) < len(b) {
		next := b[len(a)]

		if _, err := strconv.Atoi(next); err == nil {
			return -1
		}

		return +1
	}

	return 0
}

func compareForRubyGems(v, w Version) int {
	return compareRubyGemsComponents(
		strings.Split(canonicalizeRubyGemVersion(v.OriginStr), "."),
		strings.Split(canonicalizeRubyGemVersion(w.OriginStr), "."),
	)
}
