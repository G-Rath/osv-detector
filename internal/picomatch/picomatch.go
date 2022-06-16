package picomatch

import (
	"regexp"
	"strings"
)

type CompiledMatcher struct {
	Matcher *regexp.Regexp
	Exclude bool
}

func MustCompile(pattern string) CompiledMatcher {
	cm := CompiledMatcher{Exclude: strings.HasPrefix(pattern, "!")}

	pattern = strings.Trim(pattern, "!")

	pattern = regexp.QuoteMeta(pattern)
	pattern = regexp.MustCompile(`\\\*\\\*`).ReplaceAllString(pattern, `.+`)
	pattern = regexp.MustCompile(`\\\*`).ReplaceAllString(pattern, `.+(?:/|$)`)

	cm.Matcher = regexp.MustCompile("^" + pattern + "$")

	return cm
}

func (cm CompiledMatcher) Match(path string) bool {
	if cm.Matcher.MatchString(path) {
		return !cm.Exclude
	}

	return cm.Exclude
}

type CompiledMatchers []CompiledMatcher

func (cms CompiledMatchers) Matches(path string) bool {
	for _, cm := range cms {
		if cm.Match(path) {
			return true
		}
	}

	return false
}

func FromPatterns(patterns []string) CompiledMatchers {
	cms := make(CompiledMatchers, 0, len(patterns))

	for _, pattern := range patterns {
		cms = append(cms, MustCompile(pattern))
	}

	return cms
}

func Picomatch() {

}

type ToRegexOptions struct {
	Flags  []rune
	NoCase bool
}

func ToRegex(source string, options ToRegexOptions) *regexp.Regexp {
	flags := ""

	if options.Flags != nil {
		flags = string(options.Flags)
	} else if options.NoCase {
		flags = "i"
	}

	if flags != "" {
		flags = "(?" + flags + ")"
	}

	re, err := regexp.Compile(flags + source)

	if err != nil {
		return regexp.MustCompile(`/$^/`)
	}

	return re
}
