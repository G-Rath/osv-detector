package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseYarnLock_v1_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/does-not-exist")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestYarnLock_v1_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/empty.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestYarnLock_v1_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/one-package.v1.lock")

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

func TestYarnLock_v1_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/two-packages.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "concat-stream",
			Version:   "1.6.2",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "concat-map",
			Version:   "0.0.1",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v1_MultipleVersions(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/multiple-versions.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "define-properties",
			Version:   "1.1.3",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "define-property",
			Version:   "0.2.5",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "define-property",
			Version:   "1.0.0",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "define-property",
			Version:   "2.0.2",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v1_MultipleConstraints(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/multiple-constraints.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "@babel/code-frame",
			Version:   "7.12.13",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "domelementtype",
			Version:   "1.3.1",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v1_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/scoped-packages.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "@babel/code-frame",
			Version:   "7.12.11",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "@babel/compat-data",
			Version:   "7.14.0",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}

func TestYarnLock_v1_VersionsWithBuildString(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseYarnLock("fixtures/yarn/versions-with-build-strings.v1.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "css-tree",
			Version:   "1.0.0-alpha.37",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "gensync",
			Version:   "1.0.0-beta.2",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "node-fetch",
			Version:   "3.0.0-beta.9",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "resolve",
			Version:   "1.20.0",
			Ecosystem: parsers.YarnEcosystem,
		},
		{
			Name:      "resolve",
			Version:   "2.0.0-next.3",
			Ecosystem: parsers.YarnEcosystem,
		},
	})
}
