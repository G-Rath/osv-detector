package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestNpmLock_v2_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/not-json.txt")

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestNpmLock_v2_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/empty.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestNpmLock_v2_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/one-package.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "wrappy",
		Version:   "1.0.2",
		Ecosystem: parsers.NpmEcosystem,
	})
}

func TestNpmLock_v2_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/one-package-dev.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "wrappy",
		Version:   "1.0.2",
		Ecosystem: parsers.NpmEcosystem,
	})
}

func TestNpmLock_v2_TwoPackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/two-packages.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("Expected to get two packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "wrappy",
		Version:   "1.0.2",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "supports-color",
		Version:   "5.5.0",
		Ecosystem: parsers.NpmEcosystem,
	})
}

func TestNpmLock_v2_ScopedPackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/scoped-package.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("Expected to get two packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "@babel/code-frame",
		Version:   "7.0.0",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "wrappy",
		Version:   "1.0.2",
		Ecosystem: parsers.NpmEcosystem,
	})
}

func TestNpmLock_v2_NestedDependencies(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseNpmLock("fixtures/npm/nested-dependencies.v2.json")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 5 {
		t.Errorf("Expected to get five packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "postcss",
		Version:   "6.0.23",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "postcss",
		Version:   "7.0.16",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "postcss-calc",
		Version:   "7.0.1",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "supports-color",
		Version:   "6.1.0",
		Ecosystem: parsers.NpmEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "supports-color",
		Version:   "5.5.0",
		Ecosystem: parsers.NpmEcosystem,
	})
}
