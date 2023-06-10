package lockfile

import (
	"encoding/json"
	"fmt"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/filetree"
	"github.com/anchore/stereoscope/pkg/filetree/filenode"
	"github.com/anchore/stereoscope/pkg/image"
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
//
// The path is expected to be relative to within the top-level node_modules
// directory, meaning it should not start with node_modules
func isNodeModulesPackageJSON(p string) bool {
	if !strings.HasSuffix(p, "package.json") {
		return false
	}

	// todo: this should probably be moved outside this function
	p = strings.TrimPrefix(p, string(filepath.Separator))
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

func WalkNodeModulesInImage(img image.Image, pathToNodeModules string) (Lockfile, error) {
	packages := make(map[string]PackageDetails)

	err := img.SquashedTree().Walk(
		func(path file.Path, f filenode.FileNode) error {
			r, err := img.OpenPathFromSquash(path)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\n%v", err)

				return nil
			}

			var pj PackageJSON
			if err := json.NewDecoder(r).Decode(&pj); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s is not a valid JSON file: %v\n", path, err)

				return nil
			}

			packages[pj.Name+"@"+pj.Version] = PackageDetails{
				Name:      pj.Name,
				Version:   pj.Version,
				Commit:    "",
				Ecosystem: NpmEcosystem,
				CompareAs: NpmEcosystem,
			}

			return nil
		},
		&filetree.WalkConditions{
			LinkOptions: []filetree.LinkResolutionOption{},
			ShouldVisit: func(path file.Path, node filenode.FileNode) bool {
				// we only want to visit the node if:
				//   1. it is a regular file
				//   2. it is within the given node_modules directory
				//   3. it is a valid package.json that would be used by Node
				return node.FileType == file.TypeRegular &&
					strings.HasPrefix(string(path), pathToNodeModules) &&
					isNodeModulesPackageJSON(strings.TrimPrefix(string(path), pathToNodeModules))
			},
			ShouldContinueBranch: func(path file.Path, node filenode.FileNode) bool {
				// We want to avoid any symlinks as they could be cyclical, and they should
				// be safe to skip since we should end up walking their targets eventually
				return !node.IsLink()
			},
		},
	)

	return Lockfile{
		FilePath: pathToNodeModules,
		ParsedAs: "node_modules",
		Packages: pkgDetailsMapToSortedSlice(packages),
	}, err
}

func WalkNodeModules(pathToNodeModules string) (Lockfile, error) {
	packages := make(map[string]PackageDetails)

	err := filepath.Walk(pathToNodeModules, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}

		if isNodeModulesPackageJSON(strings.TrimPrefix(path, pathToNodeModules)) {
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
