package parsers_test

import (
	"fmt"
	"osv-detector/detector/parsers"
	"testing"
)

func packageToString(pkg parsers.PackageDetails) string {
	return fmt.Sprintf("%s@%s (%s)", pkg.Name, pkg.Version, pkg.Ecosystem)
}

func hasPackage(packages []parsers.PackageDetails, pkg parsers.PackageDetails) bool {
	for _, details := range packages {
		if details == pkg {
			return true
		}
	}

	return false
}

func expectPackage(t *testing.T, packages []parsers.PackageDetails, pkg parsers.PackageDetails) {
	t.Helper()

	if !hasPackage(packages, pkg) {
		t.Errorf(
			"Expected packages to include %s@%s (%s), but it did not",
			pkg.Name,
			pkg.Version,
			pkg.Ecosystem,
		)
	}
}

func findMissingPackages(actualPackages []parsers.PackageDetails, expectedPackages []parsers.PackageDetails) []parsers.PackageDetails {
	var missingPackages []parsers.PackageDetails

	for _, pkg := range actualPackages {
		if !hasPackage(expectedPackages, pkg) {
			missingPackages = append(missingPackages, pkg)
		}
	}

	return missingPackages
}

func expectPackages(t *testing.T, actualPackages []parsers.PackageDetails, expectedPackages []parsers.PackageDetails) {
	t.Helper()

	if len(expectedPackages) != len(actualPackages) {
		t.Errorf("Expected to get %d packages, but got %d", len(expectedPackages), len(actualPackages))
	}

	missingActualPackages := findMissingPackages(actualPackages, expectedPackages)
	missingExpectedPackages := findMissingPackages(expectedPackages, actualPackages)

	if len(missingActualPackages) != 0 {
		for _, unexpectedPackage := range missingActualPackages {
			t.Errorf("Did not expect %s", packageToString(unexpectedPackage))
		}
	}

	if len(missingExpectedPackages) != 0 {
		for _, unexpectedPackage := range missingExpectedPackages {
			t.Errorf("Did not find %s", packageToString(unexpectedPackage))
		}
	}
}
