package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseNpmLock_v2_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseNpmLock_v2_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseNpmLock_v2_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/empty.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseNpmLock_v2_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/one-package.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v2_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/one-package-dev.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v2_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/two-packages.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v2_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/scoped-packages.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "@babel/code-frame",
			Version:   "7.0.0",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v2_NestedDependencies(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/nested-dependencies.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "postcss",
			Version:   "6.0.23",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "postcss",
			Version:   "7.0.16",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "postcss-calc",
			Version:   "7.0.1",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "6.1.0",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}

func TestParseNpmLock_v2_NestedDependenciesDup(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/nested-dependencies-dup.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "supports-color",
			Version:   "6.1.0",
			Ecosystem: parsers.NpmEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "2.0.0",
			Ecosystem: parsers.NpmEcosystem,
		},
	})
}
