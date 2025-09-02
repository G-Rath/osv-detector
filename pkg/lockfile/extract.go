package lockfile

import (
	"cmp"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/google/osv-scalibr/converter"
	"github.com/google/osv-scalibr/extractor"
	"github.com/google/osv-scalibr/extractor/filesystem"
	scalibrfs "github.com/google/osv-scalibr/fs"
)

func extract(localPath string, extractor filesystem.Extractor, ecosystem Ecosystem) ([]PackageDetails, error) {
	info, err := os.Stat(localPath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	invs, err := extractWithExtractor(context.Background(), localPath, info, extractor)

	if err != nil {
		return nil, err
	}

	return invToPackageDetails(invs, ecosystem), nil
}

func extractWithExtractor(ctx context.Context, localPath string, info fs.FileInfo, ext filesystem.Extractor) ([]*extractor.Package, error) {
	// Create a scan input centered at the system root directory,
	// to give access to the full filesystem for each extractor.
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return nil, err
	}

	rootDir := getRootDir(absPath)
	si, err := createScanInput(absPath, rootDir, info)
	if err != nil {
		return nil, err
	}

	invs, err := ext.Extract(ctx, si)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", localPath, err)
	}

	for i := range invs.Packages {
		// Set parent extractor
		invs.Packages[i].Plugins = append(invs.Packages[i].Plugins, ext.Name())

		// Make Location relative to the scan root as we are performing local scanning
		for i2 := range invs.Packages[i].Locations {
			invs.Packages[i].Locations[i2] = filepath.Join(rootDir, invs.Packages[i].Locations[i2])
		}
	}

	slices.SortFunc(invs.Packages, inventorySort)
	invsCompact := slices.CompactFunc(invs.Packages, func(a, b *extractor.Package) bool {
		return inventorySort(a, b) == 0
	})

	return invsCompact, nil
}

func createScanInput(path string, root string, fileInfo fs.FileInfo) (*filesystem.ScanInput, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Rel will strip root from the input path.
	path, err = filepath.Rel(root, path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	si := filesystem.ScanInput{
		FS:     os.DirFS(root).(scalibrfs.FS),
		Path:   path,
		Root:   root,
		Reader: reader,
		Info:   fileInfo,
	}

	return &si, nil
}

// getRootDir returns the root directory on each system.
// On Unix systems, it'll be /
// On Windows, it will most likely be the drive (e.g. C:\)
func getRootDir(path string) string {
	if runtime.GOOS == "windows" {
		return filepath.VolumeName(path) + "\\"
	}

	if strings.HasPrefix(path, "/") {
		return "/"
	}

	return ""
}

// InventorySort is a comparator function for Inventories, to be used in
// tests with cmp.Diff to disregard the order in which the Inventories
// are reported.
func inventorySort(a, b *extractor.Package) int {
	aLoc := fmt.Sprintf("%v", a.Locations)
	bLoc := fmt.Sprintf("%v", b.Locations)

	var aExtr, bExtr string
	var aPURL, bPURL string

	aPURLStruct := converter.ToPURL(a)
	bPURLStruct := converter.ToPURL(b)

	if aPURLStruct != nil {
		aPURL = aPURLStruct.String()
	}

	if bPURLStruct != nil {
		bPURL = bPURLStruct.String()
	}

	aSourceCode := fmt.Sprintf("%v", a.SourceCode)
	bSourceCode := fmt.Sprintf("%v", b.SourceCode)

	return cmp.Or(
		cmp.Compare(aLoc, bLoc),
		cmp.Compare(a.Name, b.Name),
		cmp.Compare(a.Version, b.Version),
		cmp.Compare(aSourceCode, bSourceCode),
		cmp.Compare(aExtr, bExtr),
		cmp.Compare(aPURL, bPURL),
	)
}

func invToPackageDetails(invs []*extractor.Package, ecosystem Ecosystem) []PackageDetails {
	details := make([]PackageDetails, 0, len(invs))

	for _, inv := range invs {
		commit := ""

		if inv.SourceCode != nil {
			commit = inv.SourceCode.Commit
		}

		details = append(details, PackageDetails{
			Name:      inv.Name,
			Version:   inv.Version,
			Commit:    commit,
			Ecosystem: ecosystem,
			CompareAs: ecosystem,
		})
	}

	return details
}
