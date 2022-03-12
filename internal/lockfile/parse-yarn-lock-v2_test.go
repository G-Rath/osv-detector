package lockfile_test

import (
	"osv-detector/internal/lockfile"
	"testing"
)

func TestParseYarnLock_v2_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/does-not-exist")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestYarnLock_v2_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/empty.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestYarnLock_v2_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/one-package.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "balanced-match",
			Version:   "1.0.2",
			Ecosystem: lockfile.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/two-packages.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "compare-func",
			Version:   "2.0.0",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "concat-map",
			Version:   "0.0.1",
			Ecosystem: lockfile.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_MultipleVersions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/multiple-versions.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "debug",
			Version:   "4.3.3",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "debug",
			Version:   "2.6.9",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "debug",
			Version:   "3.2.7",
			Ecosystem: lockfile.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/scoped-packages.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@babel/cli",
			Version:   "7.16.8",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "@babel/code-frame",
			Version:   "7.16.7",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "@babel/compat-data",
			Version:   "7.16.8",
			Ecosystem: lockfile.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_VersionsWithBuildString(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseYarnLock("fixtures/yarn/versions-with-build-strings.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@nicolo-ribaudo/chokidar-2",
			Version:   "2.1.8-no-fsevents.3",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "gensync",
			Version:   "1.0.0-beta.2",
			Ecosystem: lockfile.YarnEcosystem,
		},
		{
			Name:      "eslint-plugin-jest",
			Version:   "0.0.0-use.local",
			Ecosystem: lockfile.YarnEcosystem,
		},
	})
}
