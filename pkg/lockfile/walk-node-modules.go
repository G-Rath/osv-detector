package lockfile

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// isNodeModulesPackageJSON checks if the given path is to a package.json that
// would be picked up by Node's import resolution logic; that is to say,
// it is a package.json that sits at the root of an installed module.
func isNodeModulesPackageJSON(p string) bool {
	if !strings.HasSuffix(p, "package.json") {
		return false
	}

	p = strings.TrimPrefix(p, "node_modules"+string(filepath.Separator))
	segs := strings.Split(p, "node_modules"+string(filepath.Separator))

	for _, seg := range segs {
		shouldBe := 1

		// scoped packages get installed within a dedicated directory,
		// so we expect there to be another separator to represent that
		if strings.HasPrefix(seg, "@") {
			shouldBe++
		}

		if strings.Count(seg, string(filepath.Separator)) != shouldBe {
			return false
		}
	}

	return true
}

func WalkNodeModules(pathToNodeModules string) (Lockfile, error) {
	packages := make(map[string]PackageDetails)

	err := filepath.Walk(pathToNodeModules, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}

		if isNodeModulesPackageJSON(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\n%v", err)

				return nil
			}

			var pj PackageJSON
			if err := json.Unmarshal(content, &pj); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s is not a valid JSON file: %v\n", info.Name(), err)

				return nil
			}

			packages[pj.Name+"@"+pj.Version] = PackageDetails{
				Name:      pj.Name,
				Version:   pj.Version,
				Commit:    "",
				Ecosystem: NpmEcosystem,
				CompareAs: NpmEcosystem,
			}
		}

		return nil
	})

	return Lockfile{
		FilePath: pathToNodeModules,
		ParsedAs: "node_modules",
		Packages: pkgDetailsMapToSortedSlice(packages),
	}, err
}
