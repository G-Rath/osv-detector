package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
	"github.com/g-rath/osv-detector/pkg/models"
)

func TestParsePoetryLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePoetryLock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePoetryLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePoetryLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "numpy",
			Version:   "1.23.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParsePoetryLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "proto-plus",
			Version:   "1.22.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "protobuf",
			Version:   "4.21.5",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParsePoetryLock_PackageWithMetadata(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/one-package-with-metadata.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "emoji",
			Version:   "2.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParsePoetryLock_PackageWithGitSource(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/source-git.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ike",
			Version:   "0.2.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
			Commit:    "cd66602cd29f61a2d2e7fb995fef1e61708c034d",
		},
	})
}

func TestParsePoetryLock_PackageWithLegacySource(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePoetryLock("fixtures/poetry/source-legacy.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "appdirs",
			Version:   "1.4.4",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
			Commit:    "",
		},
	})
}
