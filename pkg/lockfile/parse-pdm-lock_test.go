package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func expectNilErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
}

func TestParsePdmLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePdmLock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePdmLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/empty.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePdmLock_SinglePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/single-package.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "toml",
			Version:   "0.10.2",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
	})
}

func TestParsePdmLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/two-packages.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "toml",
			Version:   "0.10.2",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
		{
			Name:      "six",
			Version:   "1.16.0",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
	})
}

func TestParsePdmLock_PackageWithDevDependencies(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/dev-dependency.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "toml",
			Version:   "0.10.2",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
		{
			Name:      "pyroute2",
			Version:   "0.7.11",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
		{
			Name:      "win-inet-pton",
			Version:   "1.1.0",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
	})
}

func TestParsePdmLock_PackageWithOptionalDependency(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/optional-dependency.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "toml",
			Version:   "0.10.2",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
		{
			Name:      "pyroute2",
			Version:   "0.7.11",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
		{
			Name:      "win-inet-pton",
			Version:   "1.1.0",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
		},
	})
}

func TestParsePdmLock_PackageWithGitDependency(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePdmLock("testdata/pdm/git-dependency.toml")

	expectNilErr(t, err)
	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "toml",
			Version:   "0.10.2",
			Ecosystem: lockfile.PdmEcosystem,
			CompareAs: lockfile.PdmEcosystem,
			Commit:    "65bab7582ce14c55cdeec2244c65ea23039c9e6f",
		},
	})
}
