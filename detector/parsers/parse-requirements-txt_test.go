package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseRequirementsTxt_Empty(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/empty.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestParseRequirementsTxt_CommentsOnly(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/only-comments.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 0 {
		t.Errorf("Expected to get no packages, but got %d", len(packages))
	}
}

func TestParseRequirementsTxt_OneRequirementUnconstrained(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/one-requirement-unconstrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "flask",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})
}

func TestParseRequirementsTxt_OneRequirementConstrained(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/one-requirement-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("Expected to get one package, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django",
		Version:   "2.2.24",
		Ecosystem: parsers.PipEcosystem,
	})
}

func TestParseRequirementsTxt_MultipleRequirementsConstrained(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/multiple-requirements-constrained.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 13 {
		t.Errorf("Expected to get 13 packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "astroid",
		Version:   "2.5.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "beautifulsoup4",
		Version:   "4.9.3",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "boto3",
		Version:   "1.17.19",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "botocore",
		Version:   "1.20.19",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "certifi",
		Version:   "2020.12.5",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "chardet",
		Version:   "4.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "circus",
		Version:   "0.17.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "click",
		Version:   "7.1.2",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django-debug-toolbar",
		Version:   "3.2.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django-filter",
		Version:   "2.4.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django-nose",
		Version:   "1.4.7",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django-storages",
		Version:   "1.11.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "django",
		Version:   "2.2.24",
		Ecosystem: parsers.PipEcosystem,
	})
}

func TestParseRequirementsTxt_MultipleRequirementsMixed(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/multiple-requirements-mixed.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 8 {
		t.Errorf("Expected to get eight packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "flask",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "flask-cors",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "pandas",
		Version:   "0.23.4",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "numpy",
		Version:   "1.16.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "scikit-learn",
		Version:   "0.20.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "sklearn",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "requests",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "gevent",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})
}

func TestParseRequirementsTxt_FileFormatExample(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseRequirementsTxt("fixtures/pip/file-format-example.txt")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	if len(packages) != 9 {
		t.Errorf("Expected to get nine packages, but got %d", len(packages))
	}

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "pytest",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "pytest-cov",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "beautifulsoup4",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "docopt",
		Version:   "0.6.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "keyring",
		Version:   "4.1.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "coverage",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "Mopidy-Dirble",
		Version:   "1.1",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "rejected",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})

	expectPackage(t, packages, parsers.PackageDetails{
		Name:      "green",
		Version:   "0.0.0",
		Ecosystem: parsers.PipEcosystem,
	})
}
