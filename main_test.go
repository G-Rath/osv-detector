package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
	"github.com/google/go-cmp/cmp"
)

func dedent(t *testing.T, str string) string {
	t.Helper()

	// 0. replace all tabs with spaces
	str = strings.ReplaceAll(str, "\t", "  ")

	// 1. remove trailing whitespace
	re := cachedregexp.MustCompile(`\r?\n([\t ]*)$`)
	str = re.ReplaceAllString(str, "")

	// 2. if any of the lines are not indented, return as we're already dedent-ed
	re = cachedregexp.MustCompile(`(^|\r?\n)[^\t \n]`)
	if re.MatchString(str) {
		return str
	}

	// 3. find all line breaks to determine the highest common indentation level
	re = cachedregexp.MustCompile(`\n[\t ]+`)
	matches := re.FindAllString(str, -1)

	// 4. remove the common indentation from all strings
	if matches != nil {
		size := len(matches[0]) - 1

		for _, match := range matches {
			if len(match)-1 < size {
				size = len(match) - 1
			}
		}

		re := cachedregexp.MustCompile(`\n[\t ]{` + strconv.Itoa(size) + `}`)
		str = re.ReplaceAllString(str, "\n")
	}

	// 5. Remove leading whitespace.
	re = cachedregexp.MustCompile(`^\r?\n`)
	str = re.ReplaceAllString(str, "")

	return str
}

// checks if two strings are equal, treating any occurrences of `%%` in the
// expected string to mean "any text"
func areEqual(t *testing.T, actual, expect string) bool {
	t.Helper()

	expect = regexp.QuoteMeta(expect)
	expect = strings.ReplaceAll(expect, "%%", ".+")

	re := cachedregexp.MustCompile(`^` + expect + `$`)

	return re.MatchString(actual)
}

type cliTestCase struct {
	name         string
	args         []string
	wantExitCode int
	wantStdout   string
	wantStderr   string

	around func(t *testing.T) func()
}

// Attempts to normalize any file paths in the given `output` so that they can
// be compared reliably regardless of the file path separator being used.
//
// Namely, escaped forward slashes are replaced with backslashes.
func normalizeFilePaths(output string) string {
	return strings.ReplaceAll(strings.ReplaceAll(output, "\\\\", "/"), "\\", "/")
}

// wildcardDatabaseStats attempts to replace references to database stats (such as
// the number of vulnerabilities and the time that the database was last updated)
// in the output with %% wildcards, in order to reduce the noise of the cmp diff
func wildcardDatabaseStats(str string) string {
	re := cachedregexp.MustCompile(`(\w+) \(\d+ vulnerabilities, including withdrawn - last updated \w{3}, \d\d \w{3} \d{4} [012]\d:\d\d:\d\d GMT\)`)

	return re.ReplaceAllString(str, "$1 (%% vulnerabilities, including withdrawn - last updated %%)")
}

func expectAreEqual(t *testing.T, subject, actual, expect string) {
	t.Helper()

	actual = dedent(t, actual)
	expect = dedent(t, expect)

	if !areEqual(t, actual, expect) {
		if os.Getenv("TEST_NO_DIFF") == "true" {
			t.Errorf("\nactual %s does not match expected:\n got:\n%s\n\n want:\n%s", subject, actual, expect)
		} else {
			t.Errorf("\nactual %s does not match expected:\n%s", subject, cmp.Diff(expect, actual))
		}
	}
}

func testCli(t *testing.T, tc cliTestCase) {
	t.Helper()

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	ec := run(tc.args, stdoutBuffer, stderrBuffer)

	stdout := normalizeFilePaths(stdoutBuffer.String())
	stderr := normalizeFilePaths(stderrBuffer.String())

	stdout = wildcardDatabaseStats(stdout)

	if ec != tc.wantExitCode {
		t.Errorf("cli exited with code %d, not %d", ec, tc.wantExitCode)
	}

	expectAreEqual(t, "stdout output", stdout, tc.wantStdout)
	expectAreEqual(t, "stderr output", stderr, tc.wantStderr)
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name:         "",
			args:         []string{},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
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
					buildscript-gradle.lockfile
          bun.lock
					Cargo.lock
					composer.lock
					Gemfile.lock
					go.mod
					gradle.lockfile
					mix.lock
					package-lock.json
					packages.lock.json
					pdm.lock
					Pipfile.lock
					pnpm-lock.yaml
					poetry.lock
					pom.xml
					pubspec.lock
					renv.lock
					requirements.txt
					uv.lock
					yarn.lock
					csv-file
					csv-row
			`,
		},
		// only the files in the given directories are checked (no recursion)
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/")},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_EmptyDirExitCode(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// no paths should return standard error exit code
		{
			name:         "",
			args:         []string{},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// one directory without any lockfiles should result in "no lockfiles in directories" exit code
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-none")},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// two directories without any lockfiles should return "no lockfiles in directories" exit code
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-none"), filepath.FromSlash("./testdata/")},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// a path to an unknown lockfile should return standard error exit code
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-none/my-file.txt")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: `
				Error, could not determine parser for testdata/locks-none/my-file.txt
			`,
		},
		// mix and match of directory without any lockfiles and a path to an unknown lockfile should return standard exit code
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-none/my-file.txt"), filepath.FromSlash("./testdata/")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: `
				Error, could not determine parser for testdata/locks-none/my-file.txt
			`,
		},
		// when the directory does not exist, the exit code should not be for "no lockfiles in directories"
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/does/not/exist")},
			wantExitCode: 127,
			wantStdout:   "",
			// "file not found" message is different on Windows vs other OSs
			wantStderr: `
				Error reading ./testdata/does/not/exist: open ./testdata/does/not/exist: %%
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// an empty directory + a directory that does not exist should return standard exit code
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/does/not/exist"), filepath.FromSlash("./testdata/locks-none")},
			wantExitCode: 127,
			wantStdout:   "",
			// "file not found" message is different on Windows vs other OSs
			wantStderr: `
				Error reading ./testdata/does/not/exist: open ./testdata/does/not/exist: %%
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// when there are no parsable lockfiles in the directory + --json should give sensible json
		{
			name:         "",
			args:         []string{"--json", filepath.FromSlash("./testdata/locks-none")},
			wantExitCode: 128,
			wantStdout:   `{"results":[]}`,
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ListPackages(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name:         "",
			args:         []string{"--list-packages", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				testdata/locks-one/yarn.lock: found 1 package
					npm: balanced-match@1.0.2
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--list-packages", filepath.FromSlash("./testdata/locks-many")},
			wantExitCode: 0,
			wantStdout: `
				testdata/locks-many/Gemfile.lock: found 1 package
					RubyGems: ast@2.4.2
				testdata/locks-many/composer.lock: found 1 package
					Packagist: sentry/sdk@2.0.4 (4c115873c86ad5bd0ac6d962db70ca53bf8fb874)
				testdata/locks-many/yarn.lock: found 1 package
					npm: balanced-match@1.0.2
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--list-packages", filepath.FromSlash("./testdata/locks-empty")},
			wantExitCode: 0,
			wantStdout: `
				testdata/locks-empty/Gemfile.lock: found 0 packages

				testdata/locks-empty/composer.lock: found 0 packages

				testdata/locks-empty/yarn.lock: found 0 packages
			`,
			wantStderr: "",
		},
		// json results in non-json output going to stderr
		{
			name:         "",
			args:         []string{"--list-packages", "--json", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				{"results":[{"filePath":"testdata/locks-one/yarn.lock","parsedAs":"yarn.lock","packages":[{"name":"balanced-match","version":"1.0.2","ecosystem":"npm","compareAs":"npm"}]}]}
			`,
			wantStderr: `
				testdata/locks-one/yarn.lock: found 1 package
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_Lockfile(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-many")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-many/Gemfile.lock: found 1 package
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/composer.lock: found 1 package
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/locks-empty")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/locks-empty/Gemfile.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-empty/yarn.lock: found 0 packages

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// parse-as + known vulnerability exits with error code 1
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// json results in non-json output going to stderr
		{
			name:         "",
			args:         []string{"--json", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				{"results":[{"filePath":"testdata/locks-one/yarn.lock","parsedAs":"yarn.lock","packages":[{"name":"balanced-match","version":"1.0.2","ecosystem":"npm","compareAs":"npm","vulnerabilities":[],"ignored":[]}]}]}
			`,
			wantStderr: `
				Loaded the following OSV databases:
          npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_DBs(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name:         "",
			args:         []string{"--use-dbs=false", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/locks-one/yarn.lock: found 1 package

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--use-api", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					osv.dev v1 API (using batches of 1000)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-one/yarn.lock: found 1 package
					Using db osv.dev v1 API (using batches of 1000)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ParseAsSpecific(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when there is just a ":", it defaults as empty
		{
			name:         "",
			args:         []string{filepath.FromSlash(":./testdata/locks-insecure/composer.lock")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// ":" can be used as an escape (no test though because it's invalid on Windows)
		{
			name:         "",
			args:         []string{filepath.FromSlash(":./testdata/locks-insecure/my:file")},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				Error reading ./testdata/locks-insecure/my:file: open ./testdata/locks-insecure/my:file: %%
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// when a path to a file is given, parse-as is applied to that file
		{
			name:         "",
			args:         []string{filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name:         "",
			args:         []string{filepath.FromSlash("package-lock.json:./testdata/locks-insecure")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{filepath.FromSlash("package-lock.json:./testdata/locks-empty")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:


				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

			`,
			wantStderr: `
				Error, could not parse testdata/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse testdata/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{filepath.FromSlash("package-lock.json:./testdata/locks-empty"), filepath.FromSlash("package-lock.json:./testdata/locks-insecure")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)


				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found


				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: `
				Error, could not parse testdata/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse testdata/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ParseAsGlobal(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when a path to a file is given, parse-as is applied to that file
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-empty")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:


				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

			`,
			wantStderr: `
				Error, could not parse testdata/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse testdata/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-empty"), filepath.FromSlash("./testdata/locks-insecure")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)


				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found


				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: `
				Error, could not parse testdata/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse testdata/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
		// specific parse-as takes precedence over global parse-as
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", filepath.FromSlash("Gemfile.lock:./testdata/locks-empty"), filepath.FromSlash("./testdata/locks-insecure")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-empty/Gemfile.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-empty/yarn.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ParseAs_CsvRow(t *testing.T) {
	t.Parallel()

	// these tests use "--no-config" in case the repo ever has a
	// default config (which can be useful during development)
	tests := []cliTestCase{
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,,Yarp.ReverseProxy,",
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					NuGet (%% vulnerabilities, including withdrawn - last updated %%)

				-: found 1 package
					Using db NuGet (%% vulnerabilities, including withdrawn - last updated %%)

					Yarp.ReverseProxy@ is affected by the following vulnerabilities:
						GHSA-8xc6-g8xw-h2c4: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-8xc6-g8xw-h2c4)
						GHSA-jrjw-qgr2-wfcg: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-jrjw-qgr2-wfcg)

					2 known vulnerabilities found in -
			`,
			wantStderr: "",
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,,Yarp.ReverseProxy,",
				"npm,,@typescript-eslint/types,5.13.0",
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				-: found 2 packages
					Using db NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					Yarp.ReverseProxy@ is affected by the following vulnerabilities:
						GHSA-8xc6-g8xw-h2c4: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-8xc6-g8xw-h2c4)
						GHSA-jrjw-qgr2-wfcg: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-jrjw-qgr2-wfcg)

					2 known vulnerabilities found in -
			`,
			wantStderr: "",
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,,",
				"npm,,@typescript-eslint/types,5.13.0",
			},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: "Error, row 1: not enough fields (expected at least four)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ParseAs_CsvFile(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name:         "",
			args:         []string{"--parse-as", "csv-file", filepath.FromSlash("./testdata/csvs-files/two-rows.csv")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/csvs-files/two-rows.csv: found 2 packages
					Using db NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					Yarp.ReverseProxy@ is affected by the following vulnerabilities:
						GHSA-8xc6-g8xw-h2c4: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-8xc6-g8xw-h2c4)
						GHSA-jrjw-qgr2-wfcg: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-jrjw-qgr2-wfcg)

					2 known vulnerabilities found in testdata/csvs-files/two-rows.csv
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--parse-as", "csv-file", filepath.FromSlash("./testdata/csvs-files/not-a-csv.xml")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: "Error, testdata/csvs-files/not-a-csv.xml: row 1: not enough fields (expected at least four)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_Configs(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when given a path to a single lockfile, the local config should be used
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-one/yarn.lock")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/configs-one/yarn.lock: found 0 packages
					Using config at testdata/configs-one/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-two/yarn.lock")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/configs-two/yarn.lock: found 0 packages
					Using config at testdata/configs-two/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when given a path to a directory, the local config should be used for all lockfiles
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-one")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/configs-one/yarn.lock: found 0 packages
					Using config at testdata/configs-one/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-two")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/configs-two/Gemfile.lock: found 1 package
					Using config at testdata/configs-two/.osv-detector.yaml (0 ignores)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/configs-two/yarn.lock: found 0 packages
					Using config at testdata/configs-two/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// local configs should be applied based on directory of each lockfile
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-one/yarn.lock"), filepath.FromSlash("./testdata/locks-many")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/configs-one/yarn.lock: found 0 packages
					Using config at testdata/configs-one/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found

				testdata/locks-many/Gemfile.lock: found 1 package
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/composer.lock: found 1 package
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// invalid databases should be skipped
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-extra-dbs/yarn.lock")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					api#https://example.com/v1 (using batches of 1000)
					dir#file:/testdata/configs-extra-dbs (3 vulnerabilities, including withdrawn)
					zip#https://example.com/osvs/all
				testdata/configs-extra-dbs/yarn.lock: found 0 packages
					Using config at testdata/configs-extra-dbs/.osv-detector.yaml (0 ignores)
					Using db api#https://example.com/v1 (using batches of 1000)
					Using db dir#file:/testdata/configs-extra-dbs (3 vulnerabilities, including withdrawn)

					no known vulnerabilities found
			`,
			wantStderr: " failed: unable to fetch OSV database: db host returned an unexpected status code (404 Not Found)",
		},
		// databases from configs are ignored if "--no-config-databases" is passed...
		{
			name: "",
			args: []string{
				"--no-config-databases",
				filepath.FromSlash("./testdata/configs-extra-dbs/yarn.lock"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/configs-extra-dbs/yarn.lock: found 0 packages
					Using config at testdata/configs-extra-dbs/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// ...but it does still use the built-in databases
		{
			name: "",
			args: []string{
				"--config", filepath.FromSlash("./testdata/configs-extra-dbs/.osv-detector.yaml"),
				"--no-config-databases",
				filepath.FromSlash("./testdata/locks-many/yarn.lock"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-many/yarn.lock: found 1 package
					Using config at testdata/configs-extra-dbs/.osv-detector.yaml (0 ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when a global config is provided, any local configs should be ignored
		{
			name:         "",
			args:         []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-one/yarn.lock")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				testdata/configs-one/yarn.lock: found 0 packages
					Using config at testdata/my-config.yml (1 ignore)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-two")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/configs-two/Gemfile.lock: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/configs-two/yarn.lock: found 0 packages
					Using config at testdata/my-config.yml (1 ignore)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-one/yarn.lock"), filepath.FromSlash("./testdata/locks-many")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/configs-one/yarn.lock: found 0 packages
					Using config at testdata/my-config.yml (1 ignore)

					no known vulnerabilities found

				testdata/locks-many/Gemfile.lock: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/composer.lock: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				testdata/locks-many/yarn.lock: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when a local config is invalid, none of the lockfiles in that directory should
		// be checked (as the results could be different due to e.g. missing ignores)
		{
			name:         "",
			args:         []string{filepath.FromSlash("./testdata/configs-invalid"), filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)



				testdata/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: `
				Error, could not read testdata/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
				Error, could not read testdata/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
			`,
		},
		// when a global config is invalid, none of the lockfiles should be checked
		// (as the results could be different due to e.g. missing ignores)
		{
			name: "",
			args: []string{
				"--config", filepath.FromSlash("./testdata/configs-invalid/.osv-detector.yaml"),
				filepath.FromSlash("./testdata/configs-invalid"),
				filepath.FromSlash("./testdata/locks-one"),
				filepath.FromSlash("./testdata/locks-many"),
			},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				Error, could not read testdata/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_Ignores(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// no ignore count is printed if there is nothing ignored
		{
			name:         "",
			args:         []string{"--ignore", "GHSA-1234", "--ignore", "GHSA-5678", filepath.FromSlash("./testdata/locks-one")},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name: "",
			args: []string{
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// the ignored count reflects the number of vulnerabilities ignored,
		// not the number of ignores that were provided
		{
			name: "",
			args: []string{
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--ignore", "GHSA-whgm-jr23-1234",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores passed by flags are _merged_ with those specified in configs by default
		{
			name: "",
			args: []string{
				"--config", filepath.FromSlash("./testdata/my-config.yml"),
				"--ignore", "GHSA-1234",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores from configs are ignored if "--no-config-ignores" is passed
		{
			name: "",
			args: []string{
				"--no-config-ignores",
				"--config", filepath.FromSlash("./testdata/my-config.yml"),
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/my-config.yml (skipping any ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// ignores passed by flags are still respected with "--no-config-ignores"
		{
			name: "",
			args: []string{
				"--no-config-ignores",
				"--config", filepath.FromSlash("./testdata/my-config.yml"),
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/my-config.yml (skipping any ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores passed by flags are ignored with those specified in configs
		{
			name: "",
			args: []string{
				"--config", filepath.FromSlash("./testdata/my-config.yml"),
				"--ignore", "GHSA-1234",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/my-config.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func setupConfigForUpdating(t *testing.T, path string, initial string, updated string) func() {
	t.Helper()

	err := os.WriteFile(path, []byte(initial), 0600)

	if err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	return func() {
		t.Helper()

		// ensure that we always try to remove the file
		defer func() {
			if err = os.Remove(path); err != nil {
				// this will typically fail on Windows due to processes,
				// so we just treat it as a warning instead of an error
				t.Logf("could not remove test file: %v", err)
			}
		}()

		content, err := os.ReadFile(path)

		if err != nil {
			t.Fatalf("could not read test config file: %v", err)
		}

		expectAreEqual(t, "config", string(content), updated)
	}
}

func TestRun_UpdatingConfigIgnores(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when there is no existing config, nothing should be updated
		{
			name:         "",
			args:         []string{"--update-config-ignores", filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json")},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// when given an explicit config, that should be updated
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				"--config", "testdata/existing-config.yml",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/existing-config.yml (0 ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json

				Updated testdata/existing-config.yml with 1 vulnerability
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config.yml",
					"",
					`
						ignore:
  						- GHSA-whgm-jr23-g3j9
					`,
				)
			},
		},
		// when there are existing ignores
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				"--config", "testdata/existing-config-with-ignores.yml",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/existing-config-with-ignores.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)

				Updated testdata/existing-config-with-ignores.yml with 1 vulnerability
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-ignores.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
					`
						ignore:
							- GHSA-whgm-jr23-g3j9
					`,
				)
			},
		},
		// when there are existing ignores but told to ignore those
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				"--no-config-ignores",
				"--config", "testdata/existing-config-with-ignored-ignores.yml",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/existing-config-with-ignored-ignores.yml (skipping any ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure/my-package-lock.json

				Updated testdata/existing-config-with-ignored-ignores.yml with 1 vulnerability
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-ignored-ignores.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
					`
					ignore:
						- GHSA-whgm-jr23-g3j9
					`,
				)
			},
		},
		// when there are many lockfiles with one config
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				"--config", "testdata/existing-config-with-many-lockfiles.yml",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json"),
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure-many/my-package-lock.json"),
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure-nested/my-package-lock.json"),
				filepath.FromSlash("composer.lock:./testdata/locks-insecure-nested/nested/my-composer-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure/my-package-lock.json: found 1 package
					Using config at testdata/existing-config-with-many-lockfiles.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)

				testdata/locks-insecure-many/my-package-lock.json: found 6 packages
					Using config at testdata/existing-config-with-many-lockfiles.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

          ansi-regex@4.1.0 is affected by the following vulnerabilities:
            GHSA-93q8-gq69-wqmw: Inefficient Regular Expression Complexity in chalk/ansi-regex (https://github.com/advisories/GHSA-93q8-gq69-wqmw)
          nth-check@1.0.2 is affected by the following vulnerabilities:
            GHSA-rp65-9cf3-cjxr: Inefficient Regular Expression Complexity in nth-check (https://github.com/advisories/GHSA-rp65-9cf3-cjxr)
          trim-newlines@3.0.0 is affected by the following vulnerabilities:
            GHSA-7p7h-4mm5-852v: Uncontrolled Resource Consumption in trim-newlines (https://github.com/advisories/GHSA-7p7h-4mm5-852v)
          ua-parser-js@1.0.2 is affected by the following vulnerabilities:
            GHSA-fhg7-m89q-25r3: ReDoS Vulnerability in ua-parser-js version (https://github.com/advisories/GHSA-fhg7-m89q-25r3)
          word-wrap@1.2.3 is affected by the following vulnerabilities:
            GHSA-j8xg-fqg3-53r7: word-wrap vulnerable to Regular Expression Denial of Service (https://github.com/advisories/GHSA-j8xg-fqg3-53r7)

					5 known vulnerabilities found in testdata/locks-insecure-many/my-package-lock.json

				testdata/locks-insecure-nested/my-package-lock.json: found 1 package
					Using config at testdata/existing-config-with-many-lockfiles.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)

				testdata/locks-insecure-nested/nested/my-composer-lock.json: found 1 package
					Using config at testdata/existing-config-with-many-lockfiles.yml (1 ignore)
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					guzzlehttp/psr7@1.8.2 is affected by the following vulnerabilities:
						GHSA-q7rv-6hp3-vh96: Improper Input Validation in guzzlehttp/psr7 (https://github.com/advisories/GHSA-q7rv-6hp3-vh96)
						GHSA-wxmh-65f7-jcvw: Improper header name validation in guzzlehttp/psr7 (https://github.com/advisories/GHSA-wxmh-65f7-jcvw)

					2 known vulnerabilities found in testdata/locks-insecure-nested/nested/my-composer-lock.json

				Updated testdata/existing-config-with-many-lockfiles.yml with 8 vulnerabilities
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-many-lockfiles.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
					`
					ignore:
						- GHSA-7p7h-4mm5-852v
						- GHSA-93q8-gq69-wqmw
						- GHSA-fhg7-m89q-25r3
						- GHSA-j8xg-fqg3-53r7
						- GHSA-q7rv-6hp3-vh96
						- GHSA-rp65-9cf3-cjxr
						- GHSA-whgm-jr23-g3j9
						- GHSA-wxmh-65f7-jcvw
					`,
				)
			},
		},
		// when there are multiple implicit configs, it updates the right ones
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure-nested/my-package-lock.json"),
				filepath.FromSlash("composer.lock:./testdata/locks-insecure-nested/nested/my-composer-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure-nested/my-package-lock.json: found 1 package
					Using config at testdata/locks-insecure-nested/.osv-detector.yml (0 ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in testdata/locks-insecure-nested/my-package-lock.json

				testdata/locks-insecure-nested/nested/my-composer-lock.json: found 1 package
					Using config at testdata/locks-insecure-nested/nested/.osv-detector.yml (0 ignores)
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					guzzlehttp/psr7@1.8.2 is affected by the following vulnerabilities:
						GHSA-q7rv-6hp3-vh96: Improper Input Validation in guzzlehttp/psr7 (https://github.com/advisories/GHSA-q7rv-6hp3-vh96)
						GHSA-wxmh-65f7-jcvw: Improper header name validation in guzzlehttp/psr7 (https://github.com/advisories/GHSA-wxmh-65f7-jcvw)

					2 known vulnerabilities found in testdata/locks-insecure-nested/nested/my-composer-lock.json

				Updated testdata/locks-insecure-nested/.osv-detector.yml with 1 vulnerability
				Updated testdata/locks-insecure-nested/nested/.osv-detector.yml with 2 vulnerabilities
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				cleanupConfig1 := setupConfigForUpdating(t,
					"testdata/locks-insecure-nested/.osv-detector.yml",
					"ignore: []",
					`
					ignore:
						- GHSA-whgm-jr23-g3j9
					`,
				)

				cleanupConfig2 := setupConfigForUpdating(t,
					"testdata/locks-insecure-nested/nested/.osv-detector.yml",
					"ignore: []",
					`
					ignore:
						- GHSA-q7rv-6hp3-vh96
						- GHSA-wxmh-65f7-jcvw
					`,
				)

				return func() {
					cleanupConfig1()
					cleanupConfig2()
				}
			},
		},
		// when there are existing ignores, it updates them and removes patched ones
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure-many/my-package-lock.json"),
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				testdata/locks-insecure-many/my-package-lock.json: found 6 packages
					Using config at testdata/locks-insecure-many/.osv-detector.yml (3 ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

          nth-check@1.0.2 is affected by the following vulnerabilities:
            GHSA-rp65-9cf3-cjxr: Inefficient Regular Expression Complexity in nth-check (https://github.com/advisories/GHSA-rp65-9cf3-cjxr)
          ua-parser-js@1.0.2 is affected by the following vulnerabilities:
            GHSA-fhg7-m89q-25r3: ReDoS Vulnerability in ua-parser-js version (https://github.com/advisories/GHSA-fhg7-m89q-25r3)
          word-wrap@1.2.3 is affected by the following vulnerabilities:
            GHSA-j8xg-fqg3-53r7: word-wrap vulnerable to Regular Expression Denial of Service (https://github.com/advisories/GHSA-j8xg-fqg3-53r7)

					3 new vulnerabilities found in testdata/locks-insecure-many/my-package-lock.json (2 were ignored)

				Updated testdata/locks-insecure-many/.osv-detector.yml with 5 vulnerabilities
			`,
			wantStderr: "",
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/locks-insecure-many/.osv-detector.yml",
					"ignore: [GHSA-7p7h-4mm5-852v, GHSA-93q8-gq69-wqmw, GHSA-67hx-6x53-jw92]",
					`
						ignore:
							- GHSA-7p7h-4mm5-852v
							- GHSA-93q8-gq69-wqmw
							- GHSA-fhg7-m89q-25r3
							- GHSA-j8xg-fqg3-53r7
							- GHSA-rp65-9cf3-cjxr
					`,
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.around != nil {
				teardown := tt.around(t)

				defer teardown()
			}

			testCli(t, tt)
		})
	}
}

func TestRun_EndToEnd(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_ACCEPTANCE") != "true" {
		t.Skip("Skipping acceptance tests")
	}

	e2eTestdataDir := "./testdata/locks-e2e"

	files, err := os.ReadDir(e2eTestdataDir)

	if err != nil {
		t.Fatalf("%v", err)
	}

	tests := make([]cliTestCase, 0, len(files)/2)
	re := cachedregexp.MustCompile(`\d+-(.*)`)

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".out.txt") {
			continue
		}

		if f.IsDir() {
			t.Errorf("unexpected directory in e2e tests")

			continue
		}

		matches := re.FindStringSubmatch(f.Name())

		if matches == nil {
			t.Errorf("could not determine parser for %s", f.Name())

			continue
		}

		parseAs := matches[1]

		fp := filepath.FromSlash(filepath.Join(e2eTestdataDir, f.Name()))

		wantStdout, err := os.ReadFile(fp + ".out.txt")

		if err != nil {
			t.Errorf("could not read outfile for e2e fixture: %v", err)

			continue
		}

		tests = append(tests, cliTestCase{
			args:         []string{parseAs + ":" + fp},
			wantExitCode: 1,
			wantStdout:   strings.ReplaceAll(string(wantStdout), "\r", ""),
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}
