package lockfile

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type PdmLockPackage struct {
	Name     string   `toml:"name"`
	Version  string   `toml:"version"`
	Groups   []string `toml:"groups"`
	Revision string   `toml:"revision"`
}

type PdmLockFile struct {
	Version  string           `toml:"lock-version"`
	Packages []PdmLockPackage `toml:"package"`
}

const PdmEcosystem = PipEcosystem

func ParsePdmLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *PdmLockFile

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
		details := PackageDetails{
			Name:      pkg.Name,
			Version:   pkg.Version,
			Ecosystem: PdmEcosystem,
			CompareAs: PdmEcosystem,
		}

		if pkg.Revision != "" {
			details.Commit = pkg.Revision
		}

		packages = append(packages, details)
	}

	return packages, nil
}
