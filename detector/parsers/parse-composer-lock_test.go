package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func hasPackage(packages []parsers.PackageDetails, pkg parsers.PackageDetails) bool {
	for _, details := range packages {
		if details == pkg {
			return true
		}
	}

	return false
}

func expectPackage(t *testing.T, packages []parsers.PackageDetails, pkg parsers.PackageDetails) {
	if !hasPackage(packages, pkg) {
		t.Errorf(
			"Expected packages to include %s@%s (%s), but it did not",
			pkg.Name,
			pkg.Version,
			pkg.Ecosystem,
		)
	}
}

func TestParseComposerLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/not-json.txt")

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestParseComposerLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/empty.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestParseComposerLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/one-package.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "sentry/sdk",
		Version:   "2.0.4",
		Ecosystem: parsers.ComposerEcosystem,
	})
}

func TestParseComposerLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/one-package-dev.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "sentry/sdk",
		Version:   "2.0.4",
		Ecosystem: parsers.ComposerEcosystem,
	})
}

func TestParseComposerLock_TwoPackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/two-packages.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("Expected to get two packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "sentry/sdk",
		Version:   "2.0.4",
		Ecosystem: parsers.ComposerEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "theseer/tokenizer",
		Version:   "1.1.3",
		Ecosystem: parsers.ComposerEcosystem,
	})
}

func TestParseComposerLock_TwoPackageAlt(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/two-packages-alt.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("Expected to get two packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "sentry/sdk",
		Version:   "2.0.4",
		Ecosystem: parsers.ComposerEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "theseer/tokenizer",
		Version:   "1.1.3",
		Ecosystem: parsers.ComposerEcosystem,
	})
}
