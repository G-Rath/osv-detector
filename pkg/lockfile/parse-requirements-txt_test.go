package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
	"github.com/g-rath/osv-detector/pkg/models"
)

func TestParseRequirementsTxt_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/does-not-exist")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxt_Empty(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/empty.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxt_CommentsOnly(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/only-comments.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxt_OneRequirementUnconstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/one-package-unconstrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_OneRequirementConstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/one-package-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_MultipleRequirementsConstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/multiple-packages-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "astroid",
			Version:   "2.5.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "beautifulsoup4",
			Version:   "4.9.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "boto3",
			Version:   "1.17.19",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "botocore",
			Version:   "1.20.19",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "certifi",
			Version:   "2020.12.5",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "chardet",
			Version:   "4.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "circus",
			Version:   "0.17.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "click",
			Version:   "7.1.2",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django-debug-toolbar",
			Version:   "3.2.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django-filter",
			Version:   "2.4.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django-nose",
			Version:   "1.4.7",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django-storages",
			Version:   "1.11.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_MultipleRequirementsMixed(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/multiple-packages-mixed.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "flask-cors",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "numpy",
			Version:   "1.16.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "scikit-learn",
			Version:   "0.20.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "sklearn",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "requests",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "gevent",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_FileFormatExample(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/file-format-example.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "pytest",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "pytest-cov",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "beautifulsoup4",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "docopt",
			Version:   "0.6.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "keyring",
			Version:   "4.1.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "coverage",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "mopidy-dirble",
			Version:   "1.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "rejected",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "green",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_WithAddedSupport(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/with-added-support.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "twisted",
			Version:   "20.3.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_NonNormalizedNames(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/non-normalized-names.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "zope-interface",
			Version:   "5.4.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "pillow",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "twisted",
			Version:   "20.3.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_WithMultipleROptions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/with-multiple-r-options.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "flask-cors",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "numpy",
			Version:   "1.16.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "scikit-learn",
			Version:   "0.20.1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "sklearn",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "requests",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "gevent",
			Version:   "0.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_WithBadROption(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/with-bad-r-option.txt")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxt_DuplicateROptions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/duplicate-r-dev.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "django",
			Version:   "0.1.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "unittest",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_CyclicRSelf(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/cyclic-r-self.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_CyclicRComplex(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/cyclic-r-complex-1.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "cyclic-r-complex",
			Version:   "1",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "cyclic-r-complex",
			Version:   "2",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "cyclic-r-complex",
			Version:   "3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_WithPerRequirementOptions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/with-per-requirement-options.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "boto3",
			Version:   "1.26.121",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "foo",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "fooproject",
			Version:   "1.2",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "barproject",
			Version:   "1.2",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}

func TestParseRequirementsTxt_LineContinuation(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxt("fixtures/pip/line-continuation.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "foo",
			Version:   "1.2.3",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "bar",
			Version:   "4.5\\\\",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "baz",
			Version:   "7.8.9",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
		{
			Name:      "qux",
			Version:   "10.11.12",
			Ecosystem: models.EcosystemPyPI,
			CompareAs: models.EcosystemPyPI,
		},
	})
}
