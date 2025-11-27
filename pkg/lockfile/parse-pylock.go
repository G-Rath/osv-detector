package lockfile

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type pylockVCS struct {
	Type   string `toml:"type"`
	Commit string `toml:"commit-id"`
}

type pylockDirectory struct {
	Path string `toml:"path"`
}

type PylockPackage struct {
	Name      string          `toml:"name"`
	Version   string          `toml:"version"`
	VCS       pylockVCS       `toml:"vcs"`
	Directory pylockDirectory `toml:"directory"`
}

type PylockLockfile struct {
	Version  string          `toml:"lock-version"`
	Packages []PylockPackage `toml:"packages"`
}

const PylockEcosystem = PipEcosystem

func ParsePylock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *PylockLockfile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = toml.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, pkg := range parsedLockfile.Packages {
		// this is likely the root package, which is sometimes included in the lockfile
		if pkg.Version == "" && pkg.Directory.Path == "." {
			continue
		}

		pkgDetails := PackageDetails{
			Name:      pkg.Name,
			Version:   pkg.Version,
			Ecosystem: PylockEcosystem,
			CompareAs: PylockEcosystem,
		}

		if pkg.VCS.Commit != "" {
			pkgDetails.Commit = pkg.VCS.Commit
		}

		packages = append(packages, pkgDetails)
	}

	return packages, nil
}
