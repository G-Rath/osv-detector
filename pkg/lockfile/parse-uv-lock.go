package lockfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type UvLockPackageSource struct {
	Virtual string `toml:"virtual"`
	Git     string `toml:"git"`
}

type UvLockPackage struct {
	Name    string              `toml:"name"`
	Version string              `toml:"version"`
	Source  UvLockPackageSource `toml:"source"`

	// uv stores "groups" as a table under "package" after all the packages, which due
	// to how TOML works means it ends up being a property on the last package, even
	// through in this context it's a global property rather than being per-package
	Groups map[string][]UvOptionalDependency `toml:"optional-dependencies"`
}

type UvOptionalDependency struct {
	Name string `toml:"name"`
}
type UvLockFile struct {
	Version  int             `toml:"version"`
	Packages []UvLockPackage `toml:"package"`
}

const UvEcosystem = PipEcosystem

func ParseUvLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *UvLockFile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = toml.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, lockPackage := range parsedLockfile.Packages {
		// skip including the root "package", since its name and version are most likely arbitrary
		if lockPackage.Source.Virtual == "." {
			continue
		}

		_, commit, _ := strings.Cut(lockPackage.Source.Git, "#")

		packages = append(packages, PackageDetails{
			Name:      lockPackage.Name,
			Version:   lockPackage.Version,
			Ecosystem: UvEcosystem,
			CompareAs: UvEcosystem,
			Commit:    commit,
		})
	}

	return packages, nil
}
