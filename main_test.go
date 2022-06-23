package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func dedent(t *testing.T, str string) string {
	t.Helper()

	// 0. replace all tabs with spaces
	str = strings.ReplaceAll(str, "\t", "  ")

	// 1. remove trailing whitespace
	re := regexp.MustCompile(`\r?\n([\t ]*)$`)
	str = re.ReplaceAllString(str, "")

	// 2. if any of the lines are not indented, return as we're already dedent-ed
	re = regexp.MustCompile(`(^|\r?\n)[^\t \n]`)
	if re.MatchString(str) {
		return str
	}

	// 3. find all line breaks to determine the highest common indentation level
	re = regexp.MustCompile(`\n[\t ]+`)
	matches := re.FindAllString(str, -1)

	// 4. remove the common indentation from all strings
	if matches != nil {
		size := len(matches[0]) - 1

		for _, match := range matches {
			if len(match)-1 < size {
				size = len(match) - 1
			}
		}

		re := regexp.MustCompile(`\n[\t ]{` + fmt.Sprint(size) + `}`)
		str = re.ReplaceAllString(str, "\n")
	}

	// 5. Remove leading whitespace.
	re = regexp.MustCompile(`^\r?\n`)
	str = re.ReplaceAllString(str, "")

	return str
}

// checks if two strings are equal, treating any occurrences of `%%` in the
// expected string to mean "any text"
func areEqual(t *testing.T, actual, expect string) bool {
	t.Helper()

	expect = regexp.QuoteMeta(expect)
	expect = strings.ReplaceAll(expect, "%%", ".+")

	re := regexp.MustCompile(`^` + expect + `$`)

	return re.MatchString(actual)
}

func runCLI(t *testing.T, args []string) (int, string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	return run(args, stdout, stderr), stdout.String(), stderr.String()
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		{
			name:         "",
			args:         []string{},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing a lockfile (see --help for usage and flags)
			`,
		},
		{
			name:         "",
			args:         []string{"--version"},
			wantExitCode: 0,
			wantStdout:   "osv-detector dev (unknown, commit none)",
			wantStderr:   "",
		},
		{
			name:         "",
			args:         []string{"--parse-as", "my-file"},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				Don't know how to parse files as "my-file" - supported values are:
					cargo.lock
					composer.lock
					Gemfile.lock
					go.mod
					package-lock.json
					pnpm-lock.yaml
					pom.xml
					requirements.txt
					yarn.lock
					csv-file
					csv-row
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ec, stdout, stderr := runCLI(t, tt.args)

			if ec != tt.wantExitCode {
				t.Errorf("cli exited with code %d, not %d", ec, tt.wantExitCode)
			}

			if !areEqual(t, dedent(t, stdout), dedent(t, tt.wantStdout)) {
				t.Errorf("stdout\n got: \n%s\n\n want:\n%s", dedent(t, stdout), dedent(t, tt.wantStdout))
			}

			if !areEqual(t, dedent(t, stderr), dedent(t, tt.wantStderr)) {
				t.Errorf("stderr\n got:\n%s\n\n want:\n%s", dedent(t, stderr), dedent(t, tt.wantStderr))
			}
		})
	}
}
