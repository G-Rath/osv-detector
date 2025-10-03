package lockfile_test

import (
	"errors"
	"os"
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParsePnpmLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/does-not-exist")

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected \"%v\" error but got \"%v\"", os.ErrNotExist, err)
	}
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePnpmLock_InvalidYaml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/not-yaml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePnpmLock_Empty(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/empty.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePnpmLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/no-packages.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePnpmLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/one-package.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_OnePackageV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/one-package-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/one-package-dev.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_ScopedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/scoped-packages.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_ScopedPackagesV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/scoped-packages-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.57.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependencies(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/peer-dependencies.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "acorn-jsx",
			Version:   "5.3.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "acorn",
			Version:   "8.7.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependenciesV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/peer-dependencies-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "js-tokens",
			Version:   "4.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "loose-envify",
			Version:   "1.4.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react-dom",
			Version:   "18.2.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react",
			Version:   "18.2.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "scheduler",
			Version:   "0.23.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependenciesAdvanced(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/peer-dependencies-advanced.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@typescript-eslint/eslint-plugin",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/parser",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/type-utils",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/types",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/typescript-estree",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@typescript-eslint/utils",
			Version:   "5.13.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "eslint-utils",
			Version:   "3.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "eslint",
			Version:   "8.10.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "tsutils",
			Version:   "3.21.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependenciesAdvancedV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/peer-dependencies-advanced-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "js-tokens",
			Version:   "4.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "loose-envify",
			Version:   "1.4.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react-dom",
			Version:   "18.3.0-canary-ab31a9ed2-20230824",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react",
			Version:   "18.3.0-canary-ab31a9ed2-20230824",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "scheduler",
			Version:   "0.24.0-canary-ab31a9ed2-20230824",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_PeerDependenciesAdvancedRCV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/peer-dependencies-advanced-rc-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "js-tokens",
			Version:   "4.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "loose-envify",
			Version:   "1.4.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react-dom",
			Version:   "18.0.0-rc.3",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "react",
			Version:   "18.2.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "scheduler",
			Version:   "0.21.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_MultiplePackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/multiple-packages.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "aws-sdk",
			Version:   "2.1087.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "base64-js",
			Version:   "1.5.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "buffer",
			Version:   "4.9.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "events",
			Version:   "1.1.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "ieee754",
			Version:   "1.1.13",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "isarray",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "jmespath",
			Version:   "0.16.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "punycode",
			Version:   "1.3.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "querystring",
			Version:   "0.2.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "sax",
			Version:   "1.2.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "url",
			Version:   "0.10.3",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "uuid",
			Version:   "3.3.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "xml2js",
			Version:   "0.4.19",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "xmlbuilder",
			Version:   "9.0.7",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_MultipleVersions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/multiple-versions.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "uuid",
			Version:   "3.3.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "uuid",
			Version:   "8.3.2",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "xmlbuilder",
			Version:   "9.0.7",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_Tarball(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/tarball.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@my-org/my-package",
			Version:   "3.2.3",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
	})
}

func TestParsePnpmLock_Exotic(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/exotic.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "foo",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@foo/bar",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.1.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "@foo/bar",
			Version:   "1.1.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.2.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.3.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.4.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_Commits(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/commits.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "my-bitbucket-package",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "6104ae42cd32c3d724036d3964678f197b2c9cdb",
		},
		{
			Name:      "@my-scope/my-package",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "267087851ad5fac92a184749c27cd539e2fc862e",
		},
		{
			Name:      "@my-scope/my-other-package",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "fbfc962ab51eb1d754749b68c064460221fbd689",
		},
		{
			Name:      "faker-parser",
			Version:   "0.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "d2dc42a9351d4d89ec48c525e34f612b6d77993f",
		},
		{
			Name:      "mocks",
			Version:   "20.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "590f321b4eb3f692bb211bd74e22947639a6f79d",
		},
	})
}

func TestParsePnpmLock_InvalidPackagePath(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/invalid-package-path.yaml")

	expectErrContaining(t, err, "invalid dependency path")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePnpmLock_Files(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/files.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "my-file-package",
			Version:   "0.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "a-local-package",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "a-nested-local-package",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "one-up",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
		{
			Name:      "one-up-with-peer",
			Version:   "1.0.0",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
			Commit:    "",
		},
	})
}

func TestParsePnpmLock_WithAliases(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/with-aliases.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@babel/helper-validator-identifier",
			Version:   "7.22.20",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "ansi-regex",
			Version:   "6.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "ansi-regex",
			Version:   "5.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}

func TestParsePnpmLock_WithAliasesV6(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePnpmLock("testdata/pnpm/with-aliases-v6.yaml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "@babel/helper-validator-identifier",
			Version:   "7.22.20",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "ansi-regex",
			Version:   "6.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
		{
			Name:      "ansi-regex",
			Version:   "5.0.1",
			Ecosystem: lockfile.PnpmEcosystem,
			CompareAs: lockfile.PnpmEcosystem,
		},
	})
}
