package parsers_test

import (
	"errors"
	"io/ioutil"
	"osv-detector/detector/parsers"
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
		"requirements.txt",
	}

	count := 0

	for _, lockfile := range lockfiles {
		_, err := parsers.TryParse("/path/to/my/"+lockfile, "")

		if errors.Is(err, parsers.ErrParserNotFound) {
			t.Errorf("No parser was found for %s", lockfile)
		}

		count++
	}

	expectNumberOfParsersCalled(t, count)
}
