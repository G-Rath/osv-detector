package lockfile

import (
	"encoding/json"
	"errors"
	"fmt"
)

type NuGetLockPackage struct {
	Resolved string `json:"resolved"`
}

// NuGetLockfile contains the required dependency information as defined in
// https://github.com/NuGet/NuGet.Client/blob/6.5.0.136/src/NuGet.Core/NuGet.ProjectModel/ProjectLockFile/PackagesLockFileFormat.cs
type NuGetLockfile struct {
	Version      int                                    `json:"version"`
	Dependencies map[string]map[string]NuGetLockPackage `json:"dependencies"`
}

const NuGetEcosystem Ecosystem = "NuGet"

func parseNuGetLockDependencies(dependencies map[string]NuGetLockPackage) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	for name, dependency := range dependencies {
		details[name+"@"+dependency.Resolved] = PackageDetails{
			Name:      name,
			Version:   dependency.Resolved,
			Ecosystem: NuGetEcosystem,
			CompareAs: NuGetEcosystem,
		}
	}

	return details
}

func parseNuGetLock(lockfile NuGetLockfile) ([]PackageDetails, error) {
	details := map[string]PackageDetails{}

	// go through the dependencies for each framework, e.g. `net6.0` and parse
	// its dependencies, there might be different or duplicate dependencies
	// between frameworks
	for _, dependencies := range lockfile.Dependencies {
		details = mergePkgDetailsMap(details, parseNuGetLockDependencies(dependencies))
	}

	return pkgDetailsMapToSlice(details), nil
}

var ErrNuGetUnsupportedLockfileVersion = errors.New("unsupported lockfile version")

func ParseNuGetLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseNuGetLock)
}

func ParseNuGetLock(f ParsableFile) ([]PackageDetails, error) {
	var parsedLockfile *NuGetLockfile

	err := json.NewDecoder(f).Decode(&parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", err)
	}

	if parsedLockfile.Version != 1 {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", ErrNuGetUnsupportedLockfileVersion)
	}

	return parseNuGetLock(*parsedLockfile)
}
