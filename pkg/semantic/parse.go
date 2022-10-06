package semantic

import (
	"errors"
	"fmt"
	"github.com/g-rath/osv-detector/internal"
	"math/big"
	"regexp"
	"strings"
)

var ErrUnsupportedEcosystem = errors.New("unsupported ecosystem")

func MustParse(str string, ecosystem internal.Ecosystem) Version {
	v, err := Parse(str, ecosystem)

	if err != nil {
		panic(err)
	}

	return v
}

func Parse(str string, ecosystem internal.Ecosystem) (Version, error) {
	switch {
	case ecosystem == "npm":
		return parseSemverVersion(str), nil
	case ecosystem == "crates.io":
		return parseSemverVersion(str), nil
	case ecosystem == "Debian":
		return parseDebianVersion(str), nil
	case ecosystem == "RubyGems":
		return parseRubyGemsVersion(str), nil
	case ecosystem == "NuGet":
		return parseNuGetVersion(str), nil
	case ecosystem == "Packagist":
		return parsePackagistVersion(str), nil
	case ecosystem == "Go":
		return parseSemverVersion(str), nil
	case ecosystem == "Hex":
		return parseSemverVersion(str), nil
	case ecosystem == "Maven":
		return parseMavenVersion(str), nil
	case ecosystem == "PyPI":
		return parsePyPIVersion(str), nil
	case ecosystem == "Pub":
		return parseSemverVersion(str), nil
	}

	return nil, fmt.Errorf("%w %s", ErrUnsupportedEcosystem, ecosystem)
}

// SemverLikeVersion is a version that is _like_ a version as defined by the
// Semantic Version specification, except with potentially unlimited numeric
// components and a leading "v"
type SemverLikeVersion struct {
	LeadingV   bool
	Components Components
	Build      string
	Original   string
}

func ParseSemverLikeVersion(line string, maxComponents int) SemverLikeVersion {
	v := parseSemverLike(line)

	if maxComponents == -1 {
		return v
	}

	components, build := v.fetchComponentsAndBuild(maxComponents)

	return SemverLikeVersion{
		LeadingV:   v.LeadingV,
		Components: components,
		Build:      build,
		Original:   v.Original,
	}
}

func parseSemverLike(line string) SemverLikeVersion {
	var components []*big.Int
	originStr := line

	numberReg := regexp.MustCompile(`\d`)

	currentCom := ""
	foundBuild := false
	emptyComponent := false

	leadingV := strings.HasPrefix(line, "v")
	line = strings.TrimPrefix(line, "v")

	for _, c := range line {
		if foundBuild {
			currentCom += string(c)

			continue
		}

		// this is part of a component version
		if numberReg.MatchString(string(c)) {
			currentCom += string(c)

			continue
		}

		// at this point, we:
		//   1. might be parsing a component (as foundBuild != true)
		//   2. we're not looking at a part of a component (as c != number)
		//
		// so c must be either:
		//   1. a component terminator (.), or
		//   2. the start of the build string
		//
		// either way, we will be terminating the current component being
		// parsed (if any), so let's do that first
		if currentCom != "" {
			v, _ := new(big.Int).SetString(currentCom, 10)

			components = append(components, v)
			currentCom = ""

			emptyComponent = false
		}

		// a component terminator means there might be another component
		// afterwards, so don't start parsing the build string just yet
		if c == '.' {
			emptyComponent = true

			continue
		}

		// anything else is part of the build string
		foundBuild = true
		currentCom = string(c)
	}

	// if we looped over everything without finding a build string,
	// then what we were currently parsing is actually a component
	if !foundBuild && currentCom != "" {
		v, _ := new(big.Int).SetString(currentCom, 10)

		components = append(components, v)
		currentCom = ""
		emptyComponent = false
	}

	// if we ended with an empty component section,
	// prefix the build string with a '.'
	if emptyComponent {
		currentCom = "." + currentCom
	}

	// if we found no components, then the v wasn't actually leading
	if len(components) == 0 && leadingV {
		leadingV = false
		currentCom = "v" + currentCom
	}

	return SemverLikeVersion{
		LeadingV:   leadingV,
		Components: components,
		Build:      currentCom,
		Original:   originStr,
	}
}
