package parsers

import (
	"fmt"
	"path"
)

func findParser(pathToLockfile string) PackageDetailsParser {
	switch path.Base(pathToLockfile) {
	case "composer.lock":
		return ParseComposerLock
	case "package-lock.json":
		return ParseNpmLock
	default:
		return nil
	}
}

func TryParse(pathToLockfile string) ([]PackageDetails, error) {
	parser := findParser(pathToLockfile)

	if parser == nil {
		return []PackageDetails{}, fmt.Errorf("cannot parse %s", path.Base(pathToLockfile))
	}

	return parser(pathToLockfile)
}
