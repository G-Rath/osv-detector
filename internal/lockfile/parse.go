package lockfile

import (
	"errors"
	"fmt"
	"path"
	"sort"
)

func findParser(pathToLockfile string) PackageDetailsParser {
	switch pathToLockfile {
	case "cargo.lock":
		return ParseCargoLock
	case "composer.lock":
		return ParseComposerLock
	case "Gemfile.lock":
		return ParseGemfileLock
	case "package-lock.json":
		return ParseNpmLock
	case "yarn.lock":
		return ParseYarnLock
	case "go.mod":
		return ParseGoLock
	case "pnpm-lock.yaml":
		return ParsePnpmLock
	case "requirements.txt":
		return ParseRequirementsTxt
	default:
		return nil
	}
}

var ErrParserNotFound = errors.New("could not determine parser")

type Packages []PackageDetails

func toSliceOfEcosystems(ecosystemsMap map[Ecosystem]struct{}) []Ecosystem {
	ecosystems := make([]Ecosystem, 0, len(ecosystemsMap))

	for ecosystem := range ecosystemsMap {
		ecosystems = append(ecosystems, ecosystem)
	}

	return ecosystems
}

func (ps Packages) Ecosystems() []Ecosystem {
	ecosystems := make(map[Ecosystem]struct{})

	for _, pkg := range ps {
		ecosystems[pkg.Ecosystem] = struct{}{}
	}

	slicedEcosystems := toSliceOfEcosystems(ecosystems)

	sort.Slice(slicedEcosystems, func(i, j int) bool {
		return slicedEcosystems[i] < slicedEcosystems[j]
	})

	return slicedEcosystems
}

// Parse attempts to extract a collection of package details from a lockfile,
// using one of the native parsers.
//
// The parser is selected based on the name of the file, which can be overridden
// with the "parseAs" parameter.
func Parse(pathToLockfile string, parseAs string) (Packages, error) {
	if parseAs == "" {
		parseAs = path.Base(pathToLockfile)
	}

	parser := findParser(parseAs)

	if parser == nil {
		return []PackageDetails{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	packages, err := parser(pathToLockfile)

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == packages[j].Name {
			return packages[i].Version < packages[j].Version
		}

		return packages[i].Name < packages[j].Name
	})

	return packages, err
}
