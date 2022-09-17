package semantic

import (
	"regexp"
	"strconv"
	"strings"
)

func replaceOnce(pat *regexp.Regexp, src, repl string) string {
	flag := false

	return pat.ReplaceAllStringFunc(src, func(a string) string {
		if flag {
			return a
		}
		flag = true

		return pat.ReplaceAllString(a, repl)
	})
}

func canonicalizePackagistVersion(v string) string {
	v = strings.TrimPrefix(v, "v")

	v = replaceOnce(regexp.MustCompile(`[-_+]`), v, ".")
	v = replaceOnce(regexp.MustCompile(`([^\d.])(\d)`), v, "$1.$2")
	v = replaceOnce(regexp.MustCompile(`(\d)([^\d.])`), v, "$1.$2")

	return v
}

func weighPackagistBuildCharacter(str string) int {
	if strings.HasPrefix(str, "RC") {
		return 3
	}

	specials := []string{"dev", "a", "b", "rc", "#", "p"}

	for i, special := range specials {
		if strings.HasPrefix(str, special) {
			return i
		}
	}

	return 0
}

func comparePackagistSpecialVersions(a, b string) int {
	av := weighPackagistBuildCharacter(a)
	bv := weighPackagistBuildCharacter(b)

	if av > bv {
		return 1
	} else if av < bv {
		return -1
	}

	return 0
}

func comparePackagistComponents(a, b []string) int {
	min := minInt(len(a), len(b))

	var compare int

	for i := 0; i < min; i++ {
		ai, aIsNumber := convertToBigInt(a[i])
		bi, bIsNumber := convertToBigInt(b[i])

		switch {
		case aIsNumber && bIsNumber:
			compare = ai.Cmp(bi)
		case !aIsNumber && !bIsNumber:
			compare = comparePackagistSpecialVersions(a[i], b[i])
		case aIsNumber:
			compare = comparePackagistSpecialVersions("#", b[i])
		default:
			compare = comparePackagistSpecialVersions(a[i], "#")
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

		return comparePackagistComponents(a[len(b):], []string{"#"})
	}

	if len(a) < len(b) {
		next := b[len(a)]

		if _, err := strconv.Atoi(next); err == nil {
			return -1
		}

		return comparePackagistComponents([]string{"#"}, b[len(a):])
	}

	return 0
}

func compareForPackagist(a, b string) int {
	return comparePackagistComponents(
		strings.Split(canonicalizePackagistVersion(a), "."),
		strings.Split(canonicalizePackagistVersion(b), "."),
	)
}
