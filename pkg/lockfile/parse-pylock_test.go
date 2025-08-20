package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParsePylock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_NoPackages(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/empty.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_OnePackage(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/one-package.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		// ...
	})
}

func TestParsePylock_TwoPackages(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/two-packages.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		// ...
	})
}

func TestParsePylock_Example(t *testing.T) {
	t.Parallel()

	// from https://peps.python.org/pep-0751/#example
	packages, err := lockfile.ParsePylock("testdata/pylock/example.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "attrs",
			Version:   "25.1.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "cattrs",
			Version:   "24.1.2",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "numpy",
			Version:   "2.2.3",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
	})
}
