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

func TestTryParse_FindsExpectedParsers(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
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
