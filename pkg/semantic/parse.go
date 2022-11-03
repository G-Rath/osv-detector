package semantic

import (
	"errors"
	"fmt"
	"github.com/g-rath/osv-detector/internal"
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
	switch ecosystem {
	case "npm":
		return parseSemverVersion(str), nil
	case "crates.io":
		return parseSemverVersion(str), nil
	case "Debian":
		return parseDebianVersion(str), nil
	case "RubyGems":
		return parseRubyGemsVersion(str), nil
	case "NuGet":
		return parseNuGetVersion(str), nil
	case "Packagist":
		return parsePackagistVersion(str), nil
	case "Go":
		return parseSemverVersion(str), nil
	case "Hex":
		return parseSemverVersion(str), nil
	case "Maven":
		return parseMavenVersion(str), nil
	case "PyPI":
		return parsePyPIVersion(str), nil
	case "Pub":
		return parseSemverVersion(str), nil
	}

	return nil, fmt.Errorf("%w %s", ErrUnsupportedEcosystem, ecosystem)
}
