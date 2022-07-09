// nolint:testpackage // main cannot be accessed directly, so cannot use main_test
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

type cliTestCase struct {
	name         string
	args         []string
	wantExitCode int
	wantStdout   string
	wantStderr   string
}

func testCli(t *testing.T, tc cliTestCase) {
	t.Helper()

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	ec := run(tc.args, stdoutBuffer, stderrBuffer)

	stdout := stdoutBuffer.String()
	stderr := stderrBuffer.String()

	if ec != tc.wantExitCode {
		t.Errorf("cli exited with code %d, not %d", ec, tc.wantExitCode)
	}

	if !areEqual(t, dedent(t, stdout), dedent(t, tc.wantStdout)) {
		t.Errorf("stdout\n got:\n%s\n\n want:\n%s", dedent(t, stdout), dedent(t, tc.wantStdout))
	}

	if !areEqual(t, dedent(t, stderr), dedent(t, tc.wantStderr)) {
		t.Errorf("stderr\n got:\n%s\n\n want:\n%s", dedent(t, stderr), dedent(t, tc.wantStderr))
	}
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
					cargo.lock
					composer.lock
					Gemfile.lock
					go.mod
					mix.lock
					package-lock.json
					pnpm-lock.yaml
					pom.xml
					requirements.txt
					yarn.lock
					csv-file
					csv-row
			`,
		},
		// only the files in the given directories are checked (no recursion)
		{
			name:         "",
			args:         []string{"./fixtures/"},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"./fixtures/locks-none"},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// two directories without any lockfiles should return "no lockfiles in directories" exit code
		{
			name:         "",
			args:         []string{"./fixtures/locks-none", "./fixtures/"},
			wantExitCode: 128,
			wantStdout:   "",
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// a path to an unknown lockfile should return standard error exit code
		{
			name:         "",
			args:         []string{"./fixtures/locks-none/my-file.txt"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: `
				Error, could not determine parser for fixtures/locks-none/my-file.txt
			`,
		},
		// mix and match of directory without any lockfiles and a path to an unknown lockfile should return standard exit code
		{
			name:         "",
			args:         []string{"./fixtures/locks-none/my-file.txt", "./fixtures/"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: `
				Error, could not determine parser for fixtures/locks-none/my-file.txt
			`,
		},
		// when the directory does not exist, the exit code should not be for "no lockfiles in directories"
		{
			name:         "",
			args:         []string{"./fixtures/does/not/exist"},
			wantExitCode: 127,
			wantStdout:   "",
			// "file not found" message is different on Windows vs other OSs
			wantStderr: `
				Error reading ./fixtures/does/not/exist: open ./fixtures/does/not/exist: %%
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// an empty directory + a directory that does not exist should return standard exit code
		{
			name:         "",
			args:         []string{"./fixtures/does/not/exist", "./fixtures/locks-none"},
			wantExitCode: 127,
			wantStdout:   "",
			// "file not found" message is different on Windows vs other OSs
			wantStderr: `
				Error reading ./fixtures/does/not/exist: open ./fixtures/does/not/exist: %%
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
		// when there are no parsable lockfiles in the directory + --json should give sensible json
		{
			name:         "",
			args:         []string{"--json", "./fixtures/locks-none"},
			wantExitCode: 128,
			wantStdout:   `{"results":[]}`,
			wantStderr: `
				You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"--list-packages", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				fixtures/locks-one/yarn.lock: found 1 package
					npm: balanced-match@1.0.2
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--list-packages", "./fixtures/locks-many"},
			wantExitCode: 0,
			wantStdout: `
				fixtures/locks-many/Gemfile.lock: found 1 package
					RubyGems: ast@2.4.2
				fixtures/locks-many/composer.lock: found 1 package
					Packagist: sentry/sdk@2.0.4 (4c115873c86ad5bd0ac6d962db70ca53bf8fb874)
				fixtures/locks-many/yarn.lock: found 1 package
					npm: balanced-match@1.0.2
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--list-packages", "./fixtures/locks-empty"},
			wantExitCode: 0,
			wantStdout: `
				fixtures/locks-empty/Gemfile.lock: found 0 packages

				fixtures/locks-empty/composer.lock: found 0 packages

				fixtures/locks-empty/yarn.lock: found 0 packages
			`,
			wantStderr: "",
		},
		// json results in non-json output going to stderr
		{
			name:         "",
			args:         []string{"--list-packages", "--json", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				{"results":[{"filePath":"fixtures/locks-one/yarn.lock","parsedAs":"yarn.lock","packages":[{"name":"balanced-match","version":"1.0.2","ecosystem":"npm"}]}]}
			`,
			wantStderr: `
				fixtures/locks-one/yarn.lock: found 1 package
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/locks-many"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-many/Gemfile.lock: found 1 package
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"./fixtures/locks-empty"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				fixtures/locks-empty/Gemfile.lock: found 0 packages

					no known vulnerabilities found

				fixtures/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

				fixtures/locks-empty/yarn.lock: found 0 packages

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// parse-as + known vulnerability exits with error code 1
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure/my-package-lock.json"},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// json results in non-json output going to stderr
		{
			name:         "",
			args:         []string{"--json", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				{"results":[{"filePath":"fixtures/locks-one/yarn.lock","parsedAs":"yarn.lock","packages":[{"name":"balanced-match","version":"1.0.2","ecosystem":"npm","vulnerabilities":[],"ignored":[]}]}]}
			`,
			wantStderr: `
				Loaded the following OSV databases:
          npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"--use-dbs=false", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				fixtures/locks-one/yarn.lock: found 1 package

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--use-api", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					osv.dev v1 API (using batches of 1000)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
					Using db osv.dev v1 API (using batches of 1000)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func TestRun_ParseAs(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when a path to a file is given, parse-as is applied to that file
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure/my-package-lock.json"},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure"},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-empty"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:


				fixtures/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found

			`,
			wantStderr: `
				Error, could not parse fixtures/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse fixtures/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-empty", "./fixtures/locks-insecure"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)


				fixtures/locks-empty/composer.lock: found 0 packages

					no known vulnerabilities found


				fixtures/locks-insecure/composer.lock: found 0 packages

					no known vulnerabilities found

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: `
				Error, could not parse fixtures/locks-empty/Gemfile.lock: unexpected end of JSON input
				Error, could not parse fixtures/locks-empty/yarn.lock: invalid character '#' looking for beginning of value
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
				"NuGet,Yarp.ReverseProxy,",
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					NuGet (%% vulnerabilities, including withdrawn - last updated %%)

				-: found 1 package
					Using db NuGet (%% vulnerabilities, including withdrawn - last updated %%)

					Yarp.ReverseProxy@ is affected by the following vulnerabilities:
						GHSA-8xc6-g8xw-h2c4: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-8xc6-g8xw-h2c4)

					1 known vulnerability found in -
			`,
			wantStderr: "",
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,Yarp.ReverseProxy,",
				"npm,@typescript-eslint/types,5.13.0",
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

					1 known vulnerability found in -
			`,
			wantStderr: "",
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,",
				"npm,@typescript-eslint/types,5.13.0",
			},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: "Error, row 1: not enough fields (expected at least three)",
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"--parse-as", "csv-file", "./fixtures/csvs-files/two-rows.csv"},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/csvs-files/two-rows.csv: found 2 packages
					Using db NuGet (%% vulnerabilities, including withdrawn - last updated %%)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					Yarp.ReverseProxy@ is affected by the following vulnerabilities:
						GHSA-8xc6-g8xw-h2c4: YARP Denial of Service Vulnerability (https://github.com/advisories/GHSA-8xc6-g8xw-h2c4)

					1 known vulnerability found in fixtures/csvs-files/two-rows.csv
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--parse-as", "csv-file", "./fixtures/csvs-files/not-a-csv.xml"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:

			`,
			wantStderr: "Error, fixtures/csvs-files/not-a-csv.xml: row 1: not enough fields (expected at least three)",
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"./fixtures/configs-one/yarn.lock"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				fixtures/configs-one/yarn.lock: found 0 packages
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
				Loaded the following OSV databases:

				fixtures/configs-two/yarn.lock: found 0 packages
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
				Loaded the following OSV databases:

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
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-two/Gemfile.lock: found 1 package
					Using config at fixtures/configs-two/.osv-detector.yaml (0 ignores)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

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
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/configs-one/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found

				fixtures/locks-many/Gemfile.lock: found 1 package
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// invalid databases should be skipped
		{
			name:         "",
			args:         []string{"./fixtures/configs-extra-dbs/yarn.lock"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					api#https://example.com/v1 (using batches of 1000)
					dir#file:/fixtures/configs-extra-dbs (0 vulnerabilities, including withdrawn)
					zip#https://example.com/osvs/all
				fixtures/configs-extra-dbs/yarn.lock: found 0 packages
					Using config at fixtures/configs-extra-dbs/.osv-detector.yaml (0 ignores)
					Using db api#https://example.com/v1 (using batches of 1000)
					Using db dir#file:/fixtures/configs-extra-dbs (0 vulnerabilities, including withdrawn)

					no known vulnerabilities found
			`,
			wantStderr: " failed: unable to fetch OSV database: could not read OSV database archive: zip: not a valid zip file",
		},
		// databases from configs are ignored if "--no-config-databases" is passed...
		{
			name: "",
			args: []string{
				"--no-config-databases",
				"./fixtures/configs-extra-dbs/yarn.lock",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				fixtures/configs-extra-dbs/yarn.lock: found 0 packages
					Using config at fixtures/configs-extra-dbs/.osv-detector.yaml (0 ignores)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// ...but it does still use the built-in databases
		{
			name: "",
			args: []string{
				"--config", "./fixtures/configs-extra-dbs/.osv-detector.yaml",
				"--no-config-databases",
				"./fixtures/locks-many/yarn.lock",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-many/yarn.lock: found 1 package
					Using config at fixtures/configs-extra-dbs/.osv-detector.yaml (0 ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when a global config is provided, any local configs should be ignored
		{
			name:         "",
			args:         []string{"--config", "fixtures/my-config.yml", "./fixtures/configs-one/yarn.lock"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:

				fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/my-config.yml (1 ignore)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--config", "fixtures/my-config.yml", "./fixtures/configs-two"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-two/Gemfile.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/configs-two/yarn.lock: found 0 packages
					Using config at fixtures/my-config.yml (1 ignore)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		{
			name:         "",
			args:         []string{"--config", "fixtures/my-config.yml", "./fixtures/configs-one/yarn.lock", "./fixtures/locks-many"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/my-config.yml (1 ignore)

					no known vulnerabilities found

				fixtures/locks-many/Gemfile.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					Using db RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					Using db Packagist (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: "",
		},
		// when a local config is invalid, none of the lockfiles in that directory should
		// be checked (as the results could be different due to e.g. missing ignores)
		{
			name:         "",
			args:         []string{"./fixtures/configs-invalid", "./fixtures/locks-one"},
			wantExitCode: 127,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)



				fixtures/locks-one/yarn.lock: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no known vulnerabilities found
			`,
			wantStderr: `
				Error, could not read fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
				Error, could not read fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
			`,
		},
		// when a global config is invalid, none of the lockfiles should be checked
		// (as the results could be different due to e.g. missing ignores)
		{
			name: "",
			args: []string{
				"--config", "./fixtures/configs-invalid/.osv-detector.yaml",
				"./fixtures/configs-invalid",
				"./fixtures/locks-one",
				"./fixtures/locks-many",
			},
			wantExitCode: 127,
			wantStdout:   "",
			wantStderr: `
				Error, could not read fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.rawConfig
			`,
		},
	}
	for _, tt := range tests {
		tt := tt
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
			args:         []string{"--ignore", "GHSA-1234", "--ignore", "GHSA-5678", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
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
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
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
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores passed by flags are _merged_ with those specified in configs by default
		{
			name: "",
			args: []string{
				"--config", "./fixtures/my-config.yml",
				"--ignore", "GHSA-1234",
				"--parse-as", "package-lock.json",
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
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
				"--config", "./fixtures/my-config.yml",
				"--parse-as", "package-lock.json",
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 1,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using config at fixtures/my-config.yml (skipping any ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// ignores passed by flags are still respected with "--no-config-ignores"
		{
			name: "",
			args: []string{
				"--no-config-ignores",
				"--config", "./fixtures/my-config.yml",
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--parse-as", "package-lock.json",
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using config at fixtures/my-config.yml (skipping any ignores)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores passed by flags are ignored with those specified in configs
		{
			name: "",
			args: []string{
				"--config", "./fixtures/my-config.yml",
				"--ignore", "GHSA-1234",
				"--parse-as", "package-lock.json",
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loaded the following OSV databases:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					Using db npm (%% vulnerabilities, including withdrawn - last updated %%)

					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}
