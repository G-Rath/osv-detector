package lockfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/jsonc"
)

type BunLockfile struct {
	Version  int              `json:"lockfileVersion"`
	Packages map[string][]any `json:"packages"`
}

const BunEcosystem = NpmEcosystem

// structurePackageDetails returns the name, version, and commit of a package
// specified as a tuple in a bun.lock
func structurePackageDetails(a []any) (string, string, string) {
	str, ok := a[0].(string)

	if !ok {
		return "", "", ""
	}

	str, isScoped := strings.CutPrefix(str, "@")
	name, version, _ := strings.Cut(str, "@")

	if isScoped {
		name = "@" + name
	}

	version, commit, _ := strings.Cut(version, "#")

	// bun.lock does not track both the commit and version,
	// so if we have a commit then we don't have a version
	if commit != "" {
		version = ""
	}

	// file dependencies do not have a semantic version recorded
	if strings.HasPrefix(version, "file:") {
		version = ""
	}

	return name, version, commit
}

func ParseBunLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *BunLockfile

	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = json.Unmarshal(jsonc.ToJSON(lockfileContents), &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Packages))

	for _, pkg := range parsedLockfile.Packages {
		name, version, commit := structurePackageDetails(pkg)

		if name == "" && version == "" && commit == "" {
			continue
		}

		packages = append(packages, PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: BunEcosystem,
			CompareAs: BunEcosystem,
			Commit:    commit,
		})
	}

	return packages, nil
}
