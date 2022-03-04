package parsers

import (
	"encoding/json"
	"fmt"
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

func pkgDetailsMapToSlice(m map[string]PackageDetails) []PackageDetails {
	var details []PackageDetails

	for _, detail := range m {
		details = append(details, detail)
	}

	return details
}

func mergePkgDetailsMap(m1 map[string]PackageDetails, m2 map[string]PackageDetails) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	for name, detail := range m1 {
		details[name] = detail
	}

	for name, detail := range m2 {
		details[name] = detail
	}

	return details
}

func parseNpmLockDependencies(dependencies map[string]NpmLockDependency) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	for name, detail := range dependencies {
		if detail.Dependencies != nil {
			details = mergePkgDetailsMap(details, parseNpmLockDependencies(detail.Dependencies))
		}

		details[name+"@"+detail.Version] = PackageDetails{
			Name:      name,
			Version:   detail.Version,
			Ecosystem: NpmEcosystem,
		}
	}

	return details
}

func extractNpmPackageName(name string) string {
	maybeScope := path.Base(path.Dir(name))
	pkgName := path.Base(name)

	if strings.HasPrefix(maybeScope, "@") {
		pkgName = maybeScope + "/" + pkgName
	}

	return pkgName
}

func parseNpmLockPackages(packages map[string]NpmLockPackage) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	for namePath, detail := range packages {
		if namePath == "" {
			continue
		}
		finalName := extractNpmPackageName(namePath)

		details[finalName+"@"+detail.Version] = PackageDetails{
			Name:      finalName,
			Version:   detail.Version,
			Ecosystem: NpmEcosystem,
		}
	}

	return details
}

func parseNpmLock(lockfile NpmLockfile) map[string]PackageDetails {
	if lockfile.Packages != nil {
		return parseNpmLockPackages(lockfile.Packages)
	}

	return parseNpmLockDependencies(lockfile.Dependencies)
}

func ParseNpmLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *NpmLockfile

	lockfileContents, err := ioutil.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = json.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	return pkgDetailsMapToSlice(parseNpmLock(*parsedLockfile)), nil
}
