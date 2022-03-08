package parsers_test

import (
	"osv-detector/detector/parsers"
	"testing"
)

func TestParseGoLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGoLock_Invalid(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/not-go-mod.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGoLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/empty.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGoLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/one-package.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "github.com/BurntSushi/toml",
			Version:   "1.0.0",
			Ecosystem: parsers.GoEcosystem,
		},
	})
}

func TestParseGoLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/two-packages.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "github.com/BurntSushi/toml",
			Version:   "1.0.0",
			Ecosystem: parsers.GoEcosystem,
		},
		{
			Name:      "gopkg.in/yaml.v2",
			Version:   "2.4.0",
			Ecosystem: parsers.GoEcosystem,
		},
	})
}

func TestParseGoLock_IndirectPackages(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGoLock("fixtures/go/indirect-packages.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "github.com/BurntSushi/toml",
			Version:   "1.0.0",
			Ecosystem: parsers.GoEcosystem,
		},
		{
			Name:      "gopkg.in/yaml.v2",
			Version:   "2.4.0",
			Ecosystem: parsers.GoEcosystem,
		},
		{
			Name:      "github.com/mattn/go-colorable",
			Version:   "0.1.9",
			Ecosystem: parsers.GoEcosystem,
		},
		{
			Name:      "github.com/mattn/go-isatty",
			Version:   "0.0.14",
			Ecosystem: parsers.GoEcosystem,
		},
		{
			Name:      "golang.org/x/sys",
			Version:   "0.0.0-20210630005230-0f9fa26af87c",
			Ecosystem: parsers.GoEcosystem,
		},
	})
}
