package parsers

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

type NpmLockDependency struct {
	Version      string                       `json:"version"`
	Dependencies map[string]NpmLockDependency `json:"dependencies,omitempty"`
}

type NpmLockPackage struct {
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type NpmLockfile struct {
	Version int `json:"lockfileVersion"`
	// npm v1- lockfiles use "dependencies"
	Dependencies map[string]NpmLockDependency `json:"dependencies"`
	// npm v2+ lockfiles use "packages"
	Packages map[string]NpmLockPackage `json:"packages,omitempty"`
}

const NpmEcosystem Ecosystem = "npm"

func parseNpmLockDependencies(dependencies map[string]NpmLockDependency) []PackageDetails {
	var details []PackageDetails

	for name, detail := range dependencies {
		if detail.Dependencies != nil {
			details = append(details, parseNpmLockDependencies(detail.Dependencies)...)
		}

		details = append(details, PackageDetails{
			Name:      name,
			Version:   detail.Version,
			Ecosystem: NpmEcosystem,
		})
	}

	return details
}

func extractPackageName(name string) string {
	maybeScope := path.Base(path.Dir(name))
	pkgName := path.Base(name)

	if strings.HasPrefix(maybeScope, "@") {
		pkgName = maybeScope + "/" + pkgName
	}

	return pkgName
}

func parseNpmLockPackages(packages map[string]NpmLockPackage) []PackageDetails {
	var details []PackageDetails

	for namePath, detail := range packages {
		if namePath == "" {
			continue
		}

		details = append(details, PackageDetails{
			Name:      extractPackageName(namePath),
			Version:   detail.Version,
			Ecosystem: NpmEcosystem,
		})
	}

	return details
}

func parseNpmLock(lockfile NpmLockfile) []PackageDetails {
	if lockfile.Packages != nil {
		return parseNpmLockPackages(lockfile.Packages)
	}

	return parseNpmLockDependencies(lockfile.Dependencies)
}

func ParseNpmLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *NpmLockfile

	if lockfileContents, err := ioutil.ReadFile(pathToLockfile); err == nil {
		err := json.Unmarshal(lockfileContents, &parsedLockfile)

		if err != nil {
			return []PackageDetails{}, err
		}
	}

	return parseNpmLock(*parsedLockfile), nil
}
