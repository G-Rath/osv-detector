package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseNpmLockFile_v1_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLockFile_v1_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLockFile_v1_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/empty.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseNpmLockFile_v1_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/one-package.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLockFile_v1_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/one-package-dev.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLockFile_v1_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/two-packages.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLockFile_v1_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/scoped-packages.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "@babel/code-frame",
			Version:   "7.0.0",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLockFile_v1_NestedDependencies(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/nested-dependencies.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "postcss",
			Version:   "6.0.23",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "postcss",
			Version:   "7.0.16",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "postcss-calc",
			Version:   "7.0.1",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "6.1.0",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
		},
	})
}

func TestParseNpmLockFile_v1_NestedDependenciesDup(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/nested-dependencies-dup.v1.json")

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
		CompareAs: lockfile.NpmEcosystem,
	})

	expectPackage(t, packages, lockfile.PackageDetails{
		Name:      "supports-color",
		Version:   "5.5.0",
		Ecosystem: lockfile.NpmEcosystem,
		CompareAs: lockfile.NpmEcosystem,
	})

	expectPackage(t, packages, lockfile.PackageDetails{
		Name:      "supports-color",
		Version:   "2.0.0",
		Ecosystem: lockfile.NpmEcosystem,
		CompareAs: lockfile.NpmEcosystem,
	})
}

func TestParseNpmLockFile_v1_Commits(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/commits.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@segment/analytics.js-integration-facebook-pixel",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "3b1bb80b302c2e552685dc8a029797ec832ea7c9",
		},
		{
			Name:      "ansi-styles",
			Version:   "1.0.0",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "babel-preset-php",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "c5a7ba5e0ad98b8db1cb8ce105403dd4b768cced",
		},
		{
			Name:      "is-number-1",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "af885e2e890b9ef0875edd2b117305119ee5bdc5",
		},
		{
			Name:      "is-number-1",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "be5935f8d2595bcd97b05718ef1eeae08d812e10",
		},
		{
			Name:      "is-number-2",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "d5ac0584ee9ae7bd9288220a39780f155b9ad4c8",
		},
		{
			Name:      "is-number-2",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "82dcc8e914dabd9305ab9ae580709a7825e824f5",
		},
		{
			Name:      "is-number-3",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "d5ac0584ee9ae7bd9288220a39780f155b9ad4c8",
		},
		{
			Name:      "is-number-3",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "82ae8802978da40d7f1be5ad5943c9e550ab2c89",
		},
		{
			Name:      "is-number-4",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "af885e2e890b9ef0875edd2b117305119ee5bdc5",
		},
		{
			Name:      "is-number-5",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "af885e2e890b9ef0875edd2b117305119ee5bdc5",
		},
		{
			Name:      "is-number-6",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "af885e2e890b9ef0875edd2b117305119ee5bdc5",
		},
		{
			Name:      "postcss-calc",
			Version:   "7.0.1",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "raven-js",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "c2b377e7a254264fd4a1fe328e4e3cfc9e245570",
		},
		{
			Name:      "slick-carousel",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "280b560161b751ba226d50c7db1e0a14a78c2de0",
		},
	})
}

func TestParseNpmLockFile_v1_Files(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNpmLockFile("fixtures/npm/files.v1.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "lodash",
			Version:   "1.3.1",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "other_package",
			Version:   "",
			Ecosystem: lockfile.NpmEcosystem,
			CompareAs: lockfile.NpmEcosystem,
			Commit:    "",
		},
	})
}
