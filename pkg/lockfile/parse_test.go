package lockfile_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/g-rath/osv-detector/internal/reporter"
	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func expectNumberOfParsersCalled(t *testing.T, numberOfParsersCalled int) {
	t.Helper()

	directories, err := os.ReadDir(".")

	if err != nil {
		t.Fatalf("unable to read current directory: ")
	}

	count := 0

	for _, directory := range directories {
		if strings.HasPrefix(directory.Name(), "parse-") &&
			!strings.HasSuffix(directory.Name(), "_test.go") {
			count++
		}
	}

	if numberOfParsersCalled != count {
		t.Errorf(
			"Expected %d %s to have been called, but had %d",
			count,
			reporter.Form(count, "parser", "parsers"),
			numberOfParsersCalled,
		)
	}
}

func TestFindParser(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"buildscript-gradle.lockfile",
		"Cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"bun.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"gradle.lockfile",
		"mix.lock",
		"pom.xml",
		"poetry.lock",
		"uv.lock",
		"pubspec.lock",
		"Pipfile.lock",
		"requirements.txt",
	}

	for _, file := range lockfiles {
		parser, parsedAs := lockfile.FindParser(filepath.FromSlash("path/to/my/"+file), "")

		if parser == nil {
			t.Errorf("Expected a parser to be found for %s but did not", file)
		}

		if file != parsedAs {
			t.Errorf("Expected parsedAs to be %s but got %s instead", file, parsedAs)
		}
	}
}

func TestFindParser_ExplicitParseAs(t *testing.T) {
	t.Parallel()

	parser, parsedAs := lockfile.FindParser("/path/to/my/package-lock.json", "composer.lock")

	if parser == nil {
		t.Errorf("Expected a parser to be found for package-lock.json (overridden as composer.json) but did not")
	}

	if parsedAs != "composer.lock" {
		t.Errorf("Expected parsedAs to be composer.lock but got %s instead", parsedAs)
	}
}

func TestParse_FindsExpectedParsers(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"buildscript-gradle.lockfile",
		"Cargo.lock",
		"package-lock.json",
		"packages.lock.json",
		"yarn.lock",
		"bun.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"gradle.lockfile",
		"renv.lock",
		"mix.lock",
		"pom.xml",
		"pdm.lock",
		"poetry.lock",
		"uv.lock",
		"Pipfile.lock",
		"pubspec.lock",
		"requirements.txt",
	}

	count := 0

	for _, file := range lockfiles {
		_, err := lockfile.Parse("/path/to/my/"+file, "")

		if errors.Is(err, lockfile.ErrParserNotFound) {
			t.Errorf("No parser was found for %s", file)
		}

		count++
	}

	// gradle.lockfile and buildscript-gradle.lockfile use the same parser
	count--

	expectNumberOfParsersCalled(t, count)
}

func TestParse_ParserNotFound(t *testing.T) {
	t.Parallel()

	_, err := lockfile.Parse(filepath.FromSlash("/path/to/my/"), "")

	if err == nil {
		t.Errorf("Expected to get an error but did not")
	}

	if !errors.Is(err, lockfile.ErrParserNotFound) {
		t.Errorf("Did not get the expected ErrParserNotFound error - got %v instead", err)
	}
}

func TestListParsers(t *testing.T) {
	t.Parallel()

	parsers := lockfile.ListParsers()

	if first := parsers[0]; first != "buildscript-gradle.lockfile" {
		t.Errorf("Expected first element to be buildscript-gradle.lockfile, but got %s", first)
	}

	if last := parsers[len(parsers)-1]; last != "yarn.lock" {
		t.Errorf("Expected last element to be yarn.lock, but got %s", last)
	}
}

func TestLockfile_String(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  crates.io: addr2line@0.15.2",
		"  npm: @typescript-eslint/types@5.13.0",
		"  crates.io: wasi@0.10.2+wasi-snapshot-preview1",
		"  Packagist: sentry/sdk@2.0.4",
		"  crates.io: no-version",
		"  <unknown>: no-ecosystem@1.2.3",
		"  <unknown>: no-ecosystem@1.2.3 (with-commit)",
	}, "\n")

	lockf := lockfile.Lockfile{
		Packages: []lockfile.PackageDetails{
			{
				Name:      "addr2line",
				Version:   "0.15.2",
				Ecosystem: lockfile.CargoEcosystem,
			},
			{
				Name:      "@typescript-eslint/types",
				Version:   "5.13.0",
				Ecosystem: lockfile.PnpmEcosystem,
			},
			{
				Name:      "wasi",
				Version:   "0.10.2+wasi-snapshot-preview1",
				Ecosystem: lockfile.CargoEcosystem,
			},
			{
				Name:      "sentry/sdk",
				Version:   "2.0.4",
				Ecosystem: lockfile.ComposerEcosystem,
			},
			{
				Name:      "no-version",
				Version:   "",
				Ecosystem: lockfile.CargoEcosystem,
			},
			{
				Name:      "no-ecosystem",
				Version:   "1.2.3",
				Ecosystem: "",
			},
			{
				Name:      "no-ecosystem",
				Version:   "1.2.3",
				Ecosystem: "",
				Commit:    "with-commit",
			},
		},
	}

	if actual := lockf.String(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestPackages_Ecosystems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ps   lockfile.Packages
		want []lockfile.Ecosystem
	}{
		{name: "", ps: lockfile.Packages{}, want: []lockfile.Ecosystem{}},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.CargoEcosystem,
			},
		},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: lockfile.CargoEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.CargoEcosystem,
			},
		},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "@typescript-eslint/types",
					Version:   "5.13.0",
					Ecosystem: lockfile.PnpmEcosystem,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "sentry/sdk",
					Version:   "2.0.4",
					Ecosystem: lockfile.ComposerEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.ComposerEcosystem,
				lockfile.CargoEcosystem,
				lockfile.PnpmEcosystem,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.ps.Ecosystems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ecosystems() = %v, want %v", got, tt.want)
			}
		})
	}
}
