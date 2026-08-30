package lockfile

import (
	"maps"

	"github.com/google/osv-scalibr/extractor/filesystem/language/javascript/packagelockjson"
)

const NpmEcosystem Ecosystem = "npm"

func pkgDetailsMapToSlice(m map[string]PackageDetails) []PackageDetails {
	details := make([]PackageDetails, 0, len(m))

	for _, detail := range m {
		details = append(details, detail)
	}

	return details
}

func mergePkgDetailsMap(m1 map[string]PackageDetails, m2 map[string]PackageDetails) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	maps.Copy(details, m1)
	maps.Copy(details, m2)

	return details
}

func ParseNpmLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, packagelockjson.NewDefault(), NpmEcosystem)
}
