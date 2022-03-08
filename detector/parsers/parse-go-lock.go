package parsers

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"strings"
)

const GoEcosystem Ecosystem = "Go"

func ParseGoLock(pathToLockfile string) ([]PackageDetails, error) {
	lockfileContents, err := ioutil.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	parsedLockfile, err := modfile.Parse(pathToLockfile, lockfileContents, nil)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Require))

	for _, require := range parsedLockfile.Require {
		packages = append(packages, PackageDetails{
			Name:      require.Mod.Path,
			Version:   strings.TrimPrefix(require.Mod.Version, "v"),
			Ecosystem: GoEcosystem,
		})
	}

	return packages, nil
}
