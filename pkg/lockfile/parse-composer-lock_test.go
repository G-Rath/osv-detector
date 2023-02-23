package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseComposerLockFile_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLockFile_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLockFile_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/empty.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLockFile_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/one-package.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Commit:    "4c115873c86ad5bd0ac6d962db70ca53bf8fb874",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLockFile_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/one-package-dev.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Commit:    "4c115873c86ad5bd0ac6d962db70ca53bf8fb874",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLockFile_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/two-packages.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Commit:    "4c115873c86ad5bd0ac6d962db70ca53bf8fb874",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Commit:    "11336f6f84e16a720dae9d8e6ed5019efa85a0f9",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
	})
}

func TestParseComposerLockFile_TwoPackagesAlt(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLockFile("fixtures/composer/two-packages-alt.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "sentry/sdk",
			Version:   "2.0.4",
			Commit:    "4c115873c86ad5bd0ac6d962db70ca53bf8fb874",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Commit:    "11336f6f84e16a720dae9d8e6ed5019efa85a0f9",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
		},
	})
}
