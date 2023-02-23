package lockfile

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type CargoLockPackage struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type CargoLockFile struct {
	Version  int                `toml:"version"`
	Packages []CargoLockPackage `toml:"package"`
}

const CargoEcosystem Ecosystem = "crates.io"

func ParseCargoLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseCargoLock)
}

func ParseCargoLock(r io.Reader) ([]PackageDetails, error) {
	var parsedLockfile *CargoLockFile

	_, err := toml.NewDecoder(r).Decode(&parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, lockPackage := range parsedLockfile.Packages {
		packages = append(packages, PackageDetails{
			Name:      lockPackage.Name,
			Version:   lockPackage.Version,
			Ecosystem: CargoEcosystem,
			CompareAs: CargoEcosystem,
		})
	}

	return packages, nil
}
