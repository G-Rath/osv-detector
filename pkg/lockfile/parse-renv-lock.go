package lockfile

import (
	"encoding/json"
	"fmt"
	"os"
)

type RenvPackage struct {
	Package    string `json:"Package"`
	Version    string `json:"Version"`
	Repository string `json:"Repository"`
}

type RenvLockfile struct {
	Packages map[string]RenvPackage `json:"Packages"`
}

const CRANEcosystem Ecosystem = "CRAN"

func ParseRenvLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *RenvLockfile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = json.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, pkg := range parsedLockfile.Packages {
		// currently we assume that unless a package is explicitly for a different
		// repository, it is a CRAN package (even if its Source is not Repository)
		if pkg.Repository != "" && pkg.Repository != string(CRANEcosystem) {
			continue
		}

		packages = append(packages, PackageDetails{
			Name:      pkg.Package,
			Version:   pkg.Version,
			Ecosystem: CRANEcosystem,
			CompareAs: CRANEcosystem,
		})
	}

	return packages, nil
}
