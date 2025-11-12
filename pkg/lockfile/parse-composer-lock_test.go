package lockfile_test

import (
	"errors"
	"os"
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseComposerLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/does-not-exist")

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected \"%v\" error but got \"%v\"", os.ErrNotExist, err)
	}
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/empty.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseComposerLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/one-package.json")

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

func TestParseComposerLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/one-package-dev.json")

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

func TestParseComposerLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/two-packages.json")

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

func TestParseComposerLock_TwoPackagesAlt(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/two-packages-alt.json")

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

func TestParseComposerLock_DrupalPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseComposerLock("testdata/composer/drupal-packages.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "drupal/core",
			Version:   "10.4.5",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
			Commit:    "5247dbaa65b42b601058555f4a8b2bd541f5611f",
		},
		{
			Name:      "drupal/tfa",
			Version:   "2.0.0-alpha4",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
			Commit:    "",
		},
		{
			Name:      "drupal/field_time",
			Version:   "1.0.0-beta5",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
			Commit:    "",
		},
		{
			Name:      "theseer/tokenizer",
			Version:   "1.1.3",
			Ecosystem: lockfile.ComposerEcosystem,
			CompareAs: lockfile.ComposerEcosystem,
			Commit:    "11336f6f84e16a720dae9d8e6ed5019efa85a0f9",
		},
	})
}
