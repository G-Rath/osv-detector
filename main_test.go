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

func TestRun_ListPackages(t *testing.T) {
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
			args:         []string{"--list-packages", "./fixtures/locks-one"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

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
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

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
				Loading OSV databases for the following ecosystems:

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
				Loading OSV databases for the following ecosystems:
          npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
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

func TestRun_Lockfile(t *testing.T) {
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
		// parse-as + known vulnerability exits with error code 1
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure/my-package-lock.json"},
			wantExitCode: 1,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/locks-insecure/my-package-lock.json: found 1 package
					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in ./fixtures/locks-insecure/my-package-lock.json
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
				Loading OSV databases for the following ecosystems:
          npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-one/yarn.lock: found 1 package
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

func TestRun_ParseAs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		// when a path to a file is given, parse-as is applied to that file
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure/my-package-lock.json"},
			wantExitCode: 1,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/locks-insecure/my-package-lock.json: found 1 package
					ansi-html@0.0.1 is affected by the following vulnerabilities:
						GHSA-whgm-jr23-g3j9: Uncontrolled Resource Consumption in ansi-html (https://github.com/advisories/GHSA-whgm-jr23-g3j9)

					1 known vulnerability found in ./fixtures/locks-insecure/my-package-lock.json
			`,
			wantStderr: "",
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name:         "",
			args:         []string{"--parse-as", "package-lock.json", "./fixtures/locks-insecure"},
			wantExitCode: 1,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/locks-insecure/composer.lock: found 0 packages
					no known vulnerabilities found

				fixtures/locks-insecure/my-package-lock.json: found 1 package
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
				Loading OSV databases for the following ecosystems:


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
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)


				fixtures/locks-empty/composer.lock: found 0 packages
					no known vulnerabilities found


				fixtures/locks-insecure/composer.lock: found 0 packages
					no known vulnerabilities found

				fixtures/locks-insecure/my-package-lock.json: found 1 package
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
		// when a global config is provided, any local configs should be ignored
		{
			name:         "",
			args:         []string{"--config", "fixtures/my-config.yml", "./fixtures/configs-one/yarn.lock"},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:

				./fixtures/configs-one/yarn.lock: found 0 packages
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
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)

				fixtures/configs-two/Gemfile.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
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
				Loading OSV databases for the following ecosystems:
					RubyGems (%% vulnerabilities, including withdrawn - last updated %%)
					Packagist (%% vulnerabilities, including withdrawn - last updated %%)
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/configs-one/yarn.lock: found 0 packages
					Using config at fixtures/my-config.yml (1 ignore)
					no known vulnerabilities found

				fixtures/locks-many/Gemfile.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					no known vulnerabilities found

				fixtures/locks-many/composer.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
					no known vulnerabilities found

				fixtures/locks-many/yarn.lock: found 1 package
					Using config at fixtures/my-config.yml (1 ignore)
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
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)



				fixtures/locks-one/yarn.lock: found 1 package
					no known vulnerabilities found
			`,
			wantStderr: `
				Error, could not read fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.Config
				Error, could not read fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.Config
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
			wantStdout: "",
			wantStderr: `
				Error, could not read ./fixtures/configs-invalid/.osv-detector.yaml: yaml: unmarshal errors:
					line 1: cannot unmarshal !!str ` + "`ignore ...`" + ` into configer.Config
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

func TestRun_Ignores(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		// no ignore count is printed if there is nothing ignored
		{
			name:         "",
			args:         []string{"--ignore", "GHSA-1234", "--ignore", "GHSA-5678", "./fixtures/locks-one"},
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
			name: "",
			args: []string{
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--parse-as", "package-lock.json",
				"./fixtures/locks-insecure/my-package-lock.json",
			},
			wantExitCode: 0,
			wantStdout: `
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/locks-insecure/my-package-lock.json: found 1 package
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
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/locks-insecure/my-package-lock.json: found 1 package
					no new vulnerabilities found (1 was ignored)
			`,
			wantStderr: "",
		},
		// ignores passed by flags are _merged_ with those specified in configs
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
				Loading OSV databases for the following ecosystems:
					npm (%% vulnerabilities, including withdrawn - last updated %%)

				./fixtures/locks-insecure/my-package-lock.json: found 1 package
					Using config at ./fixtures/my-config.yml (1 ignore)
					no new vulnerabilities found (1 was ignored)
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
