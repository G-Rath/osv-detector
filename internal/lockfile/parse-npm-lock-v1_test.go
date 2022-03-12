package lockfile_test

import (
	"osv-detector/internal/lockfile"
	"testing"
)

func TestParseNpmLock_v1_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLock_v1_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLock_v1_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/empty.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLock_v1_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/one-package.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v1_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/one-package-dev.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v1_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/two-packages.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v1_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/scoped-packages.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "@babel/code-frame",
			Version:   "7.0.0",
			Ecosystem: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v1_NestedDependencies(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/nested-dependencies.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "postcss",
			Version:   "6.0.23",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "postcss",
			Version:   "7.0.16",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "postcss-calc",
			Version:   "7.0.1",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "6.1.0",
			Ecosystem: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v1_NestedDependenciesDup(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLock("fixtures/npm/nested-dependencies-dup.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	// todo: convert to using expectPackages w/ listing all expected packages
	if len(packages) != 39 {
		t.Errorf("Expected to get two packages, but got %d", len(packages))
	}

	expectPackage(t, packages, lockfile.PackageDetails{
		Name:      "supports-color",
		Version:   "6.1.0",
		Ecosystem: lockfile.NpmEcosystem,
	})

	expectPackage(t, packages, lockfile.PackageDetails{
		Name:      "supports-color",
		Version:   "5.5.0",
		Ecosystem: lockfile.NpmEcosystem,
	})

	expectPackage(t, packages, lockfile.PackageDetails{
		Name:      "supports-color",
		Version:   "2.0.0",
		Ecosystem: lockfile.NpmEcosystem,
	})
}
