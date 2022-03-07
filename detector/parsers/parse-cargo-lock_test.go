package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseCargoLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseCargoLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseCargoLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseCargoLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "addr2line",
			Version:   "0.15.2",
			Ecosystem: parsers.CargoEcosystem,
		},
	})
}

func TestParseCargoLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "addr2line",
			Version:   "0.15.2",
			Ecosystem: parsers.CargoEcosystem,
		},
		{
			Name:      "syn",
			Version:   "1.0.73",
			Ecosystem: parsers.CargoEcosystem,
		},
	})
}

func TestParseCargoLock_PackageWithBuildString(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseCargoLock("fixtures/cargo/package-with-build-string.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "wasi",
			Version:   "0.10.2+wasi-snapshot-preview1",
			Ecosystem: parsers.CargoEcosystem,
		},
	})
}
