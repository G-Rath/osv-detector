package lockfile

import (
	"encoding/json"
	"fmt"
)

type PipenvPackage struct {
	Version string `json:"version"`
}

type PipenvLock struct {
	Packages    map[string]PipenvPackage `json:"default"`
	PackagesDev map[string]PipenvPackage `json:"develop"`
}

const PipenvEcosystem = PipEcosystem

func ParsePipenvLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParsePipenvLock)
}

func ParsePipenvLock(f ParsableFile) ([]PackageDetails, error) {
	var parsedLockfile *PipenvLock

	err := json.NewDecoder(f).Decode(&parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", err)
	}

	packages := make(map[string]PackageDetails)

	for name, pipenvPackage := range parsedLockfile.Packages {
		if pipenvPackage.Version == "" {
			continue
		}

		version := pipenvPackage.Version[2:]

		packages[name+"@"+version] = PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: PipenvEcosystem,
			CompareAs: PipenvEcosystem,
		}
	}

	for name, pipenvPackage := range parsedLockfile.PackagesDev {
		if pipenvPackage.Version == "" {
			continue
		}

		version := pipenvPackage.Version[2:]

		packages[name+"@"+version] = PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: PipenvEcosystem,
			CompareAs: PipenvEcosystem,
		}
	}

	return pkgDetailsMapToSlice(packages), nil
}
