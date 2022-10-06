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
