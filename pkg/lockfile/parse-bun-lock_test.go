package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseBunLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseBunLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/not-json.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseBunLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/empty.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseBunLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/one-package.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/one-package-dev.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_OnePackageBadTuple(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/bad-tuple.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/two-packages.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "has-flag",
			Version:   "4.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_SamePackageDifferentGroups(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/same-package-different-groups.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "has-flag",
			Version:   "3.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/scoped-packages.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.62.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_ScopedPackagesMixed(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/scoped-packages-mixed.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@babel/code-frame",
			Version:   "7.26.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "@babel/helper-validator-identifier",
			Version:   "7.25.9",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "js-tokens",
			Version:   "4.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "picocolors",
			Version:   "1.1.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "wrappy",
			Version:   "1.0.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_OptionalPackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/optional-package.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.14.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "fsevents",
			Version:   "0.3.8",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "nan",
			Version:   "2.22.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_PeerDependenciesImplicit(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/peer-dependencies-implicit.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn-jsx",
			Version:   "5.3.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "acorn",
			Version:   "8.14.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_PeerDependenciesExplicit(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/peer-dependencies-explicit.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn-jsx",
			Version:   "5.3.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "acorn",
			Version:   "8.14.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_NestedDependenciesDups(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/nested-dependencies-dup.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ansi-styles",
			Version:   "4.3.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "chalk",
			Version:   "4.1.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "color-convert",
			Version:   "2.0.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "color-name",
			Version:   "1.1.4",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "2.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "7.2.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "4.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_NestedDependencies(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/nested-dependencies.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ansi-styles",
			Version:   "4.3.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "chalk",
			Version:   "4.1.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "color-convert",
			Version:   "2.0.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "color-name",
			Version:   "1.1.4",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "2.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "5.5.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "7.2.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "3.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "4.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_Aliases(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/alias.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "has-flag",
			Version:   "4.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "7.2.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "supports-color",
			Version:   "6.1.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
		{
			Name:      "has-flag",
			Version:   "3.0.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
		},
	})
}

func TestParseBunLock_Commits(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/commits.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@babel/helper-plugin-utils",
			Version:   "7.26.5",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "@babel/helper-string-parser",
			Version:   "7.25.9",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "@babel/helper-validator-identifier",
			Version:   "7.25.9",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "@babel/parser",
			Version:   "7.26.5",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "@babel/types",
			Version:   "7.26.5",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "@prettier/sync",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "527e8ce",
		},
		{
			Name:      "babel-preset-php",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "1ae6dc1267500360b411ec711b8aeac8c68b2246",
		},
		{
			Name:      "is-number",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "98e8ff1",
		},
		{
			Name:      "is-number",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "d5ac058",
		},
		{
			Name:      "is-number",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "b7aef34",
		},
		{
			Name:      "jquery",
			Version:   "3.7.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "lodash",
			Version:   "1.3.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "make-synchronized",
			Version:   "0.2.9",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "php-parser",
			Version:   "2.2.0",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "prettier",
			Version:   "3.4.2",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "raven-js",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "91ef2d4",
		},
		{
			Name:      "slick-carousel",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "fc6f7d8",
		},
		{
			Name:      "stopwords",
			Version:   "0.0.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
	})
}

func TestParseBunLock_Files(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseBunLock("fixtures/bun/files.json5")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "etag",
			Version:   "",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
		{
			Name:      "lodash",
			Version:   "1.3.1",
			Ecosystem: lockfile.BunEcosystem,
			CompareAs: lockfile.BunEcosystem,
			Commit:    "",
		},
	})
}
