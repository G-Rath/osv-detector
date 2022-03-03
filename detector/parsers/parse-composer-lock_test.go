package parsers_test

import (
	"osv-detector/detector/parsers"
	"strings"
	"testing"
)

func TestParseComposerLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseComposerLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/not-json.txt")

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	if !strings.Contains(err.Error(), "could not parse") {
		t.Errorf("Expected to get \"could not parse\" error, but got \"%v\"", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseComposerLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/empty.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseComposerLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/one-package.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: parsers.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/one-package-dev.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: parsers.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_TwoPackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/two-packages.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: parsers.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Ecosystem: parsers.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_TwoPackageAlt(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseComposerLock("fixtures/composer/two-packages-alt.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: parsers.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Ecosystem: parsers.ComposerEcosystem,
		},
	})
}
