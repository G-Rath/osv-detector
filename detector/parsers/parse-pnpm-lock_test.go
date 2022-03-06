package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParsePnpmLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParsePnpmLock_InvalidYaml(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/not-yaml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParsePnpmLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/empty.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParsePnpmLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/one-package.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/one-package-dev.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/scoped-packages.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependencies(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/peer-dependencies.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "acorn-jsx",
			Version:   "5.3.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependenciesAdvanced(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/peer-dependencies-advanced.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "@typescript-eslint/eslint-plugin",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/parser",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/type-utils",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/typescript-estree",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/utils",
			Version:   "5.13.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "eslint-utils",
			Version:   "3.0.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "eslint",
			Version:   "8.10.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "tsutils",
			Version:   "3.21.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_MultiplePackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/multiple-packages.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "aws-sdk",
			Version:   "2.1087.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "base64-js",
			Version:   "1.5.1",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "buffer",
			Version:   "4.9.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "events",
			Version:   "1.1.1",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "ieee754",
			Version:   "1.1.13",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "isarray",
			Version:   "1.0.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "jmespath",
			Version:   "0.16.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "punycode",
			Version:   "1.3.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "querystring",
			Version:   "0.2.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "sax",
			Version:   "1.2.1",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "url",
			Version:   "0.10.3",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "uuid",
			Version:   "3.3.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "xml2js",
			Version:   "0.4.19",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "xmlbuilder",
			Version:   "9.0.7",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_MultipleVersions(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/multiple-versions.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "uuid",
			Version:   "3.3.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "uuid",
			Version:   "8.3.2",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "xmlbuilder",
			Version:   "9.0.7",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_Exotic(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParsePnpmLock("fixtures/pnpm/exotic.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "foo",
			Version:   "1.0.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@foo/bar",
			Version:   "1.0.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.1.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "@foo/bar",
			Version:   "1.1.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.2.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.3.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.4.0",
			Ecosystem: parsers.PnpmEcosystem,
		},
	})
}
