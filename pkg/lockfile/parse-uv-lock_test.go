package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseUvLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseUvLock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseUvLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseUvLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "emoji",
			Version:   "2.14.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
	})
}

func TestParseUvLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "emoji",
			Version:   "2.14.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "protobuf",
			Version:   "4.25.5",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
	})
}

func TestParseUvLock_SourceGit(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/source-git.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ruff",
			Version:   "0.8.1",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
			Commit:    "84748be16341b76e073d117329f7f5f4ee2941ad",
		},
	})
}

func TestParseUvLock_GroupedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseUvLock("fixtures/uv/grouped-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "emoji",
			Version:   "2.14.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "click",
			Version:   "8.1.7",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "colorama",
			Version:   "0.4.6",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "black",
			Version:   "24.10.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "flake8",
			Version:   "7.1.1",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "mccabe",
			Version:   "0.7.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "mypy-extensions",
			Version:   "1.0.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "packaging",
			Version:   "24.2",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "pathspec",
			Version:   "0.12.1",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "platformdirs",
			Version:   "4.3.6",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "pycodestyle",
			Version:   "2.12.1",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "pyflakes",
			Version:   "3.2.0",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "tomli",
			Version:   "2.2.1",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
		{
			Name:      "typing-extensions",
			Version:   "4.12.2",
			Ecosystem: lockfile.UvEcosystem,
			CompareAs: lockfile.UvEcosystem,
		},
	})
}
