package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseRenvLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRenvLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRenvLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRenvLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "morning",
			Version:   "0.1.0",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
	})
}

func TestParseRenvLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "markdown",
			Version:   "1.0",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
		{
			Name:      "mime",
			Version:   "0.7",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
	})
}

func TestParseRenvLock_WithMixedSources(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/with-mixed-sources.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "markdown",
			Version:   "1.0",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
		{
			Name:      "mime",
			Version:   "0.12.1",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
	})
}
func TestParseRenvLock_WithoutRepository(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRenvLock("fixtures/renv/without-repository.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "morning",
			Version:   "0.1.0",
			Ecosystem: lockfile.CRANEcosystem,
			CompareAs: lockfile.CRANEcosystem,
		},
	})
}
