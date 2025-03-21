package lockfile

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
	"gopkg.in/yaml.v3"
)

var errInvalidPackagePath = errors.New("invalid package path")

type PnpmLockPackageResolution struct {
	Tarball string `yaml:"tarball"`
	Commit  string `yaml:"commit"`
	Repo    string `yaml:"repo"`
	Type    string `yaml:"type"`
}

type PnpmLockPackage struct {
	Resolution PnpmLockPackageResolution `yaml:"resolution"`
	Name       string                    `yaml:"name"`
	Version    string                    `yaml:"version"`
}

type PnpmLockfile struct {
	Version  float64                    `yaml:"lockfileVersion"`
	Packages map[string]PnpmLockPackage `yaml:"packages,omitempty"`
}

type pnpmLockfileV6 struct {
	Version  string                     `yaml:"lockfileVersion"`
	Packages map[string]PnpmLockPackage `yaml:"packages,omitempty"`
}

func (l *PnpmLockfile) UnmarshalYAML(value *yaml.Node) error {
	var lockfileV6 pnpmLockfileV6

	if err := value.Decode(&lockfileV6); err != nil {
		return fmt.Errorf("%w", err)
	}

	parsedVersion, err := strconv.ParseFloat(lockfileV6.Version, 64)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	l.Version = parsedVersion
	l.Packages = lockfileV6.Packages

	return nil
}

const PnpmEcosystem = NpmEcosystem

func startsWithNumber(str string) bool {
	matcher := cachedregexp.MustCompile(`^\d`)

	return matcher.MatchString(str)
}

// extractPnpmPackageNameAndVersion parses a dependency path, attempting to
// extract the name and version of the package it represents
func extractPnpmPackageNameAndVersion(dependencyPath string, lockfileVersion float64) (string, string, error) {
	// file dependencies must always have a name property to be installed,
	// and their dependency path never has the version encoded, so we can
	// skip trying to extract either from their dependency path
	if strings.HasPrefix(dependencyPath, "file:") {
		return "", "", nil
	}

	// v9.0 specifies the dependencies as <package>@<version> rather than as a path
	if lockfileVersion == 9.0 {
		dependencyPath = strings.Trim(dependencyPath, "'")
		dependencyPath, isScoped := strings.CutPrefix(dependencyPath, "@")

		name, version, _ := strings.Cut(dependencyPath, "@")

		if isScoped {
			name = "@" + name
		}

		return name, version, nil
	}

	parts := strings.Split(dependencyPath, "/")

	if len(parts) == 1 {
		return "", "", errInvalidPackagePath
	}

	var name string

	parts = parts[1:]

	if strings.HasPrefix(parts[0], "@") {
		name = strings.Join(parts[:2], "/")
		parts = parts[2:]
	} else {
		name = parts[0]
		parts = parts[1:]
	}

	version := ""

	if len(parts) != 0 {
		version = parts[0]
	}

	if version == "" {
		name, version = parseNameAtVersion(name)
	}

	if version == "" || !startsWithNumber(version) {
		return "", "", nil
	}

	// peer dependencies in v5 lockfiles are attached to the end of the version
	// with an "_", so we always want the first element if an "_" is present
	version, _, _ = strings.Cut(version, "_")

	return name, version, nil
}

func parseNameAtVersion(value string) (name string, version string) {
	// look for pattern "name@version", where name is allowed to contain zero or more "@"
	matches := cachedregexp.MustCompile(`^(.+)@([\w.-]+)(?:\(|$)`).FindStringSubmatch(value)

	if len(matches) != 3 {
		return name, ""
	}

	return matches[1], matches[2]
}

func parsePnpmLock(lockfile PnpmLockfile) ([]PackageDetails, error) {
	packages := make([]PackageDetails, 0, len(lockfile.Packages))

	for s, pkg := range lockfile.Packages {
		name, version, err := extractPnpmPackageNameAndVersion(s, lockfile.Version)

		if err != nil {
			return nil, err
		}

		// "name" is only present if it's not in the dependency path and takes
		// priority over whatever name we think we've extracted (if any)
		if pkg.Name != "" {
			name = pkg.Name
		}

		// "version" is only present if it's not in the dependency path and takes
		// priority over whatever version we think we've extracted (if any)
		if pkg.Version != "" {
			version = pkg.Version
		}

		if name == "" || version == "" {
			continue
		}

		commit := pkg.Resolution.Commit

		if strings.HasPrefix(pkg.Resolution.Tarball, "https://codeload.github.com") {
			re := cachedregexp.MustCompile(`https://codeload\.github\.com(?:/[\w-.]+){2}/tar\.gz/(\w+)$`)
			matched := re.FindStringSubmatch(pkg.Resolution.Tarball)

			if matched != nil {
				commit = matched[1]
			}
		}

		packages = append(packages, PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: PnpmEcosystem,
			CompareAs: PnpmEcosystem,
			Commit:    commit,
		})
	}

	return packages, nil
}

func ParsePnpmLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *PnpmLockfile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = yaml.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	// this will happen if the file is empty
	if parsedLockfile == nil {
		parsedLockfile = &PnpmLockfile{}
	}

	packageDetails, err := parsePnpmLock(*parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	return packageDetails, nil
}
