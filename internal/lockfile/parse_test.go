package lockfile_test

import (
	"errors"
	"io/ioutil"
	"osv-detector/internal/lockfile"
	"strings"
	"testing"
)

func expectNumberOfParsersCalled(t *testing.T, numberOfParsersCalled int) {
	t.Helper()

	directories, err := ioutil.ReadDir(".")

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
			"Expected %d parsers to have been called, but had %d",
			count,
			numberOfParsersCalled,
		)
	}
}

func TestFindParser(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"pom.xml",
		"requirements.txt",
	}

	for _, file := range lockfiles {
		parser, parsedAs := lockfile.FindParser("/path/to/my/"+file, "")

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
		"cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"pom.xml",
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

	expectNumberOfParsersCalled(t, count)
}

func TestParse_ParserNotFound(t *testing.T) {
	t.Parallel()

	_, err := lockfile.Parse("/path/to/my/", "")

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

	if first := parsers[0]; first != "cargo.lock" {
		t.Errorf("Expected first element to be cargo.lock, but got %s", first)
	}

	if last := parsers[len(parsers)-1]; last != "yarn.lock" {
		t.Errorf("Expected last element to be requirements.txt, but got %s", last)
	}
}

func TestLockfile_ToString(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  crates.io: addr2line@0.15.2",
		"  npm: @typescript-eslint/types@5.13.0",
		"  crates.io: wasi@0.10.2+wasi-snapshot-preview1",
		"  Packagist: sentry/sdk@2.0.4",
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
		},
	}

	if actual := lockf.ToString(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}
