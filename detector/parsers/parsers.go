package parsers

import (
	"fmt"
	"path"
)

func findParser(pathToLockfile string) PackageDetailsParser {
	switch pathToLockfile {
	case "composer.lock":
		return ParseComposerLock
	case "package-lock.json":
		return ParseNpmLock
	default:
		return nil
	}
}

func TryParse(pathToLockfile string, parseAs string) ([]PackageDetails, error) {
	if parseAs == "" {
		parseAs = path.Base(pathToLockfile)
	}

	parser := findParser(parseAs)

	if parser == nil {
		return []PackageDetails{}, fmt.Errorf("cannot parse %s", path.Base(pathToLockfile))
	}

	return parser(pathToLockfile)
}
