package lockfile

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type PylockPackage struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type PylockLockfile struct {
	Packages []PylockPackage `toml:"packages"`
}

const PylockEcosystem = PipEcosystem

func ParsePylock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *PylockLockfile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = toml.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, pkg := range parsedLockfile.Packages {
		packages = append(packages, PackageDetails{
			Name:      pkg.Name,
			Version:   pkg.Version,
			Ecosystem: PylockEcosystem,
			CompareAs: PylockEcosystem,
		})
	}

	return packages, nil
}
