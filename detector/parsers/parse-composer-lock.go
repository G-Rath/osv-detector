package parsers

import (
	"encoding/json"
	"io/ioutil"
)

type ComposerPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ComposerLock struct {
	Packages []ComposerPackage `json:"packages"`
}

const ComposerEcosystem Ecosystem = "Packagist"

func ParseComposerLock(pathToLockfile string) ([]PackageDetails, error) {
	var packages []PackageDetails
	var parsedLockfile *ComposerLock

	if lockfileContents, err := ioutil.ReadFile(pathToLockfile); err == nil {
		err := json.Unmarshal(lockfileContents, &parsedLockfile)

		if err != nil {
			return packages, err
		}
	}

	for _, composerPackage := range parsedLockfile.Packages {
		packages = append(packages, PackageDetails{
			Name:      composerPackage.Name,
			Version:   composerPackage.Version,
			Ecosystem: ComposerEcosystem,
		})
	}

	return packages, nil
}
