package lockfile_test

import (
	"errors"
	"os"
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseCargoLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/does-not-exist")

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected \"%v\" error but got \"%v\"", os.ErrNotExist, err)
	}
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseCargoLock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseCargoLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseCargoLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "addr2line",
			Version:   "0.15.2",
			Ecosystem: lockfile.CargoEcosystem,
			CompareAs: lockfile.CargoEcosystem,
		},
	})
}

func TestParseCargoLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "addr2line",
			Version:   "0.15.2",
			Ecosystem: lockfile.CargoEcosystem,
			CompareAs: lockfile.CargoEcosystem,
		},
		{
			Name:      "syn",
			Version:   "1.0.73",
			Ecosystem: lockfile.CargoEcosystem,
			CompareAs: lockfile.CargoEcosystem,
		},
	})
}

func TestParseCargoLock_PackageWithBuildString(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseCargoLock("testdata/cargo/package-with-build-string.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wasi",
			Version:   "0.10.2+wasi-snapshot-preview1",
			Ecosystem: lockfile.CargoEcosystem,
			CompareAs: lockfile.CargoEcosystem,
		},
	})
}
