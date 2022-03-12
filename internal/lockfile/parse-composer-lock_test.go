package lockfile_test

import (
	"osv-detector/internal/lockfile"
	"testing"
)

func TestParseComposerLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/empty.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/one-package.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/one-package-dev.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/two-packages.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: lockfile.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Ecosystem: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLock_TwoPackagesAlt(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("fixtures/composer/two-packages-alt.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Ecosystem: lockfile.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Ecosystem: lockfile.ComposerEcosystem,
		},
	})
}
