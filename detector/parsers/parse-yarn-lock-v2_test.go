package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseYarnLock_v2_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/does-not-exist")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestYarnLock_v2_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/empty.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestYarnLock_v2_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/one-package.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "balanced-match",
			Version:   "1.0.2",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/two-packages.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "compare-func",
			Version:   "2.0.0",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "concat-map",
			Version:   "0.0.1",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_MultipleVersions(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/multiple-versions.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "debug",
			Version:   "4.3.3",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "debug",
			Version:   "2.6.9",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "debug",
			Version:   "3.2.7",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v2_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/scoped-packages.v2.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "@babel/cli",
			Version:   "7.16.8",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "@babel/code-frame",
			Version:   "7.16.7",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "@babel/compat-data",
			Version:   "7.16.8",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}
