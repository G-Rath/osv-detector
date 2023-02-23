package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseRequirementsTxtFile_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxtFile_Empty(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/empty.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxtFile_CommentsOnly(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/only-comments.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxtFile_OneRequirementUnconstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/one-package-unconstrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxtFile_OneRequirementConstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/one-package-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxtFile_MultipleRequirementsConstrained(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/multiple-packages-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "astroid",
			Version:   "2.5.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "beautifulsoup4",
			Version:   "4.9.3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "boto3",
			Version:   "1.17.19",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "botocore",
			Version:   "1.20.19",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "certifi",
			Version:   "2020.12.5",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "chardet",
			Version:   "4.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "circus",
			Version:   "0.17.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "click",
			Version:   "7.1.2",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django-debug-toolbar",
			Version:   "3.2.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django-filter",
			Version:   "2.4.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django-nose",
			Version:   "1.4.7",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django-storages",
			Version:   "1.11.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxtFile_MultipleRequirementsMixed(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/multiple-packages-mixed.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "flask-cors",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "numpy",
			Version:   "1.16.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "scikit-learn",
			Version:   "0.20.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "sklearn",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "requests",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "gevent",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxtFile_FileFormatExample(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/file-format-example.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "pytest",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "pytest-cov",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "beautifulsoup4",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "docopt",
			Version:   "0.6.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "keyring",
			Version:   "4.1.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "coverage",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "mopidy-dirble",
			Version:   "1.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "rejected",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "green",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		// todo: requires -r support
		// {
		// 	Name:      "django",
		// 	Version:   "2.2.24",
		// 	Ecosystem: lockfile.PipEcosystem,
		// 	CompareAs: lockfile.PipEcosystem,
		// },
	})
}

func TestParseRequirementsTxtFile_WithAddedSupport(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/with-added-support.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "twisted",
			Version:   "20.3.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxtFile_NonNormalizedNames(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/non-normalized-names.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "zope-interface",
			Version:   "5.4.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "pillow",
			Version:   "1.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "twisted",
			Version:   "20.3.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_WithMultipleROptions(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/with-multiple-r-options.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flask",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "flask-cors",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "numpy",
			Version:   "1.16.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "scikit-learn",
			Version:   "0.20.1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "sklearn",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "requests",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "gevent",
			Version:   "0.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "django",
			Version:   "2.2.24",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_WithBadROption(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/with-bad-r-option.txt")

	expectErrContaining(t, err, "could not open")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseRequirementsTxt_DuplicateROptions(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/duplicate-r-dev.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "django",
			Version:   "0.1.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "unittest",
			Version:   "1.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_CyclicRSelf(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/cyclic-r-self.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "pandas",
			Version:   "0.23.4",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "requests",
			Version:   "1.2.3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_CyclicRComplex(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/cyclic-r-complex-1.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "cyclic-r-complex",
			Version:   "1",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "cyclic-r-complex",
			Version:   "2",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "cyclic-r-complex",
			Version:   "3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_WithPerRequirementOptions(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/with-per-requirement-options.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "boto3",
			Version:   "1.26.121",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "foo",
			Version:   "1.0.0",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "fooproject",
			Version:   "1.2",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "barproject",
			Version:   "1.2",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}

func TestParseRequirementsTxt_LineContinuation(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseRequirementsTxtFile("fixtures/pip/line-continuation.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "foo",
			Version:   "1.2.3",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "bar",
			Version:   "4.5\\\\",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "baz",
			Version:   "7.8.9",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
		{
			Name:      "qux",
			Version:   "10.11.12",
			Ecosystem: lockfile.PipEcosystem,
			CompareAs: lockfile.PipEcosystem,
		},
	})
}
