package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParsePylock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_InvalidToml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/not-toml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_NoPackages(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/empty.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_OnePackage(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/one-package.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		// ...
	})
}

func TestParsePylock_TwoPackages(t *testing.T) {
	t.Skip("todo: need a fixture")

	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/two-packages.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		// ...
	})
}

func TestParsePylock_Example(t *testing.T) {
	t.Parallel()

	// from https://peps.python.org/pep-0751/#example
	packages, err := lockfile.ParsePylock("testdata/pylock/example.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "attrs",
			Version:   "25.1.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "cattrs",
			Version:   "24.1.2",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "numpy",
			Version:   "2.2.3",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
	})
}

func TestParsePylock_PackageWithCommits(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/commits.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "click",
			Version:   "8.2.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "mleroc",
			Version:   "0.1.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
			Commit:    "735093f03c4d8be70bfaaae44074ac92d7419b6d",
		},
		{
			Name:      "packaging",
			Version:   "24.2",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "pathspec",
			Version:   "0.12.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "python-dateutil",
			Version:   "2.9.0.post0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "scikit-learn",
			Version:   "1.6.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "tqdm",
			Version:   "4.67.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
	})
}

func TestParsePylock_CreatedByPipWithJustSelf(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePylock("testdata/pylock/pip-just-self.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePylock_CreatedByPip(t *testing.T) {
	t.Parallel()

	// from https://peps.python.org/pep-0751/#example
	packages, err := lockfile.ParsePylock("testdata/pylock/pip-full.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "annotated-types",
			Version:   "0.7.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "packaging",
			Version:   "25.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "pyproject-toml",
			Version:   "0.1.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "setuptools",
			Version:   "80.9.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "wheel",
			Version:   "0.45.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
	})
}

func TestParsePylock_CreatedByPdm(t *testing.T) {
	t.Parallel()

	// from https://peps.python.org/pep-0751/#example
	packages, err := lockfile.ParsePylock("testdata/pylock/pdm-full.toml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "certifi",
			Version:   "2025.1.31",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "chardet",
			Version:   "3.0.4",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "charset-normalizer",
			Version:   "2.0.12",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "colorama",
			Version:   "0.3.9",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "idna",
			Version:   "2.7",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "py",
			Version:   "1.4.34",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "pytest",
			Version:   "3.2.5",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "requests",
			Version:   "2.27.1",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "setuptools",
			Version:   "39.2.0",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
		{
			Name:      "urllib3",
			Version:   "1.26.20",
			Ecosystem: lockfile.PylockEcosystem,
			CompareAs: lockfile.PylockEcosystem,
		},
	})
}
