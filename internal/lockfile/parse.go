package lockfile

import (
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
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

type Lockfile struct {
	FilePath string   `json:"filePath"`
	ParsedAs string   `json:"parsedAs"`
	Packages Packages `json:"packages"`
}

func (l Lockfile) ToString() string {
	lines := make([]string, 0, len(l.Packages))

	for _, details := range l.Packages {
		lines = append(lines,
			fmt.Sprintf("  %s: %s@%s", details.Ecosystem, details.Name, details.Version),
		)
	}

	return strings.Join(lines, "\n")
}

// Parse attempts to extract a collection of package details from a lockfile,
// using one of the native parsers.
//
// The parser is selected based on the name of the file, which can be overridden
// with the "parseAs" parameter.
func Parse(pathToLockfile string, parseAs string) (Lockfile, error) {
	if parseAs == "" {
		parseAs = path.Base(pathToLockfile)
	}

	parser := findParser(parseAs)

	if parser == nil {
		return Lockfile{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	packages, err := parser(pathToLockfile)

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == packages[j].Name {
			return packages[i].Version < packages[j].Version
		}

		return packages[i].Name < packages[j].Name
	})

	return Lockfile{
		FilePath: pathToLockfile,
		ParsedAs: parseAs,
		Packages: packages,
	}, err
}
