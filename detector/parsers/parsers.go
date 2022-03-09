package parsers

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

func TryParse(pathToLockfile string, parseAs string) ([]PackageDetails, error) {
	if parseAs == "" {
		parseAs = path.Base(pathToLockfile)
	}

	parser := findParser(parseAs)

	if parser == nil {
		return []PackageDetails{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	packages, err := parser(pathToLockfile)

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages, err
}
