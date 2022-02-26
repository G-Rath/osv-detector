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

func ParseComposerLock(pathToLockfile string) (EcosystemPackages, error) {
	ecosystemPackages := EcosystemPackages{
		Ecosystem: "Packagist",
	}
	var parsedLockfile *ComposerLock

	if lockfileContents, err := ioutil.ReadFile(pathToLockfile); err == nil {
		err := json.Unmarshal(lockfileContents, &parsedLockfile)

		if err != nil {
			return ecosystemPackages, err
		}
	}

	for _, composerPackage := range parsedLockfile.Packages {
		ecosystemPackages.Packages = append(ecosystemPackages.Packages, EcosystemPackage{
			Name:    composerPackage.Name,
			Version: composerPackage.Version,
		})
	}

	return ecosystemPackages, nil
}
