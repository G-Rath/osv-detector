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
		{
			name:         "",
			args:         []string{"./fixtures/locks-none"},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing a lockfile (see --help for usage and flags)
			`,
		},
		{
			name:         "",
			args:         []string{"./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/locks-many"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-many/Gemfile.lock: found 1 package
					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/locks-empty"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:

				fixtures/locks-empty/Gemfile.lock: found 0 packages
					no known vulnerabilities found

				fixtures/locks-empty/composer.lock: found 0 packages
					no known vulnerabilities found

				fixtures/locks-empty/yarn.lock: found 0 packages
					no known vulnerabilities found
			`,
			wantStderr: "",
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

func TestRun_Configs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		// when given a path to a single lockfile, the local config should be used
		{
			name:         "",
			args:         []string{"./fixtures/configs-one/yarn.lock"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:

				./fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/configs-one/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/configs-two/yarn.lock"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:

				./fixtures/configs-two/yarn.lock: found 0 packages
					Using config at fixtures/configs-two/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when given a path to a directory, the local config should be used for all lockfiles
		{
			name:         "",
			args:         []string{"./fixtures/configs-one"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:

				fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/configs-one/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/configs-two"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-two/Gemfile.lock: found 1 package
					Using config at fixtures/configs-two/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found

				fixtures/configs-two/yarn.lock: found 0 packages
					Using config at fixtures/configs-two/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// local configs should be applied based on directory of each lockfile
		{
			name:         "",
			args:         []string{"./fixtures/configs-one/yarn.lock", "./fixtures/locks-many"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/configs-one/.osv-detector.yaml (0 ignores)
					no known vulnerabilities found

				fixtures/locks-many/Gemfile.lock: found 1 package
					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					no known vulnerabilities found
			`,
			wantStderr: "",
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