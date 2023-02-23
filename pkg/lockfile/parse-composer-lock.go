package lockfile

import (
	"encoding/json"
	"fmt"
	"io"
)

type ComposerPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    struct {
		Reference string `json:"reference"`
	} `json:"dist"`
}

type ComposerLock struct {
	Packages    []ComposerPackage `json:"packages"`
	PackagesDev []ComposerPackage `json:"packages-dev"`
}

const ComposerEcosystem Ecosystem = "Packagist"

func ParseComposerLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseComposerLock)
}

func ParseComposerLock(r io.Reader) ([]PackageDetails, error) {
	var parsedLockfile *ComposerLock

	err := json.NewDecoder(r).Decode(&parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", err)
	}

	packages := make(
		[]PackageDetails,
		0,
		// len cannot return negative numbers, but the types can't reflect that
		uint64(len(parsedLockfile.Packages))+uint64(len(parsedLockfile.PackagesDev)),
	)

	for _, composerPackage := range parsedLockfile.Packages {
		packages = append(packages, PackageDetails{
			Name:      composerPackage.Name,
			Version:   composerPackage.Version,
			Commit:    composerPackage.Dist.Reference,
			Ecosystem: ComposerEcosystem,
			CompareAs: ComposerEcosystem,
		})
	}

	for _, composerPackage := range parsedLockfile.PackagesDev {
		packages = append(packages, PackageDetails{
			Name:      composerPackage.Name,
			Version:   composerPackage.Version,
			Commit:    composerPackage.Dist.Reference,
			Ecosystem: ComposerEcosystem,
			CompareAs: ComposerEcosystem,
		})
	}

	return packages, nil
}
