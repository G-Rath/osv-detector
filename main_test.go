package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	m.Run()

	dirty, err := snaps.Clean(m, snaps.CleanOpts{Sort: true})

	if err != nil {
		//nolint:forbidigo
		fmt.Println("Error cleaning snaps:", err)
		os.Exit(1)
	}
	if dirty {
		//nolint:forbidigo
		fmt.Println("Some snapshots were outdated.")
		os.Exit(1)
	}
}

type cliTestCase struct {
	name string
	args []string
	exit int

	around func(t *testing.T) func()
}

func testCli(t *testing.T, tc cliTestCase) {
	t.Helper()

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	ec := run(tc.args, stdoutBuffer, stderrBuffer)

	stdout := normalizeSnapshot(t, stdoutBuffer.String())
	stderr := normalizeSnapshot(t, stderrBuffer.String())

	if ec != tc.exit {
		t.Errorf("cli exited with code %d, not %d", ec, tc.exit)
	}

	snaps.MatchSnapshot(t, stdout)
	snaps.MatchSnapshot(t, stderr)
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		{
			name: "",
			args: []string{},
			exit: 127,
		},
		{
			name: "",
			args: []string{"--version"},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--parse-as", "my-file"},
			exit: 127,
		},
		// only the files in the given directories are checked (no recursion)
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/")},
			exit: 128,
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
			name: "",
			args: []string{},
			exit: 127,
		},
		// one directory without any lockfiles should result in "no lockfiles in directories" exit code
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-none")},
			exit: 128,
		},
		// two directories without any lockfiles should return "no lockfiles in directories" exit code
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-none"), filepath.FromSlash("./testdata/")},
			exit: 128,
		},
		// a path to an unknown lockfile should return standard error exit code
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-none/my-file.txt")},
			exit: 127,
		},
		// mix and match of directory without any lockfiles and a path to an unknown lockfile should return standard exit code
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-none/my-file.txt"), filepath.FromSlash("./testdata/")},
			exit: 127,
		},
		// when the directory does not exist, the exit code should not be for "no lockfiles in directories"
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/does/not/exist")},
			exit: 127,
			// "file not found" message is different on Windows vs other OSs
		},
		// an empty directory + a directory that does not exist should return standard exit code
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/does/not/exist"), filepath.FromSlash("./testdata/locks-none")},
			exit: 127,
			// "file not found" message is different on Windows vs other OSs
		},
		// when there are no parsable lockfiles in the directory + --json should give sensible json
		{
			name: "",
			args: []string{"--json", filepath.FromSlash("./testdata/locks-none")},
			exit: 128,
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
			name: "",
			args: []string{"--list-packages", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--list-packages", filepath.FromSlash("./testdata/locks-many")},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--list-packages", filepath.FromSlash("./testdata/locks-empty")},
			exit: 0,
		},
		// json results in non-json output going to stderr
		{
			name: "",
			args: []string{"--list-packages", "--json", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
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
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
		},
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-many")},
			exit: 0,
		},
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/locks-empty")},
			exit: 0,
		},
		// parse-as + known vulnerability exits with error code 1
		{
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json")},
			exit: 1,
		},
		// json results in non-json output going to stderr
		{
			name: "",
			args: []string{"--json", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
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
			name: "",
			args: []string{"--use-dbs=false", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--use-api", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
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
			name: "",
			args: []string{filepath.FromSlash(":./testdata/locks-insecure/composer.lock")},
			exit: 0,
		},
		// ":" can be used as an escape (no test though because it's invalid on Windows)
		{
			name: "",
			args: []string{filepath.FromSlash(":./testdata/locks-insecure/my:file")},
			exit: 127,
		},
		// when a path to a file is given, parse-as is applied to that file
		{
			name: "",
			args: []string{filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json")},
			exit: 1,
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name: "",
			args: []string{filepath.FromSlash("package-lock.json:./testdata/locks-insecure")},
			exit: 1,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name: "",
			args: []string{filepath.FromSlash("package-lock.json:./testdata/locks-empty")},
			exit: 127,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name: "",
			args: []string{filepath.FromSlash("package-lock.json:./testdata/locks-empty"), filepath.FromSlash("package-lock.json:./testdata/locks-insecure")},
			exit: 127,
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
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json")},
			exit: 1,
		},
		// when a path to a directory is given, parse-as is applied to all files in the directory
		{
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-insecure")},
			exit: 1,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-empty")},
			exit: 127,
		},
		// files that error on parsing don't stop parsable files from being checked
		{
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("./testdata/locks-empty"), filepath.FromSlash("./testdata/locks-insecure")},
			exit: 127,
		},
		// specific parse-as takes precedence over global parse-as
		{
			name: "",
			args: []string{"--parse-as", "package-lock.json", filepath.FromSlash("Gemfile.lock:./testdata/locks-empty"), filepath.FromSlash("./testdata/locks-insecure")},
			exit: 1,
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
			exit: 1,
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,,Yarp.ReverseProxy,",
				"npm,,@typescript-eslint/types,5.13.0",
			},
			exit: 1,
		},
		{
			name: "",
			args: []string{
				"--no-config",
				"--parse-as", "csv-row",
				"NuGet,,",
				"npm,,@typescript-eslint/types,5.13.0",
			},
			exit: 127,
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
			name: "",
			args: []string{"--parse-as", "csv-file", filepath.FromSlash("./testdata/csvs-files/two-rows.csv")},
			exit: 1,
		},
		{
			name: "",
			args: []string{"--parse-as", "csv-file", filepath.FromSlash("./testdata/csvs-files/not-a-csv.xml")},
			exit: 127,
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
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-one/yarn.lock")},
			exit: 0,
		},
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-two/yarn.lock")},
			exit: 0,
		},
		// when given a path to a directory, the local config should be used for all lockfiles
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-one")},
			exit: 0,
		},
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-two")},
			exit: 0,
		},
		// local configs should be applied based on directory of each lockfile
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-one/yarn.lock"), filepath.FromSlash("./testdata/locks-many")},
			exit: 0,
		},
		// invalid databases should be skipped
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-extra-dbs/yarn.lock")},
			exit: 127,
		},
		// databases from configs are ignored if "--no-config-databases" is passed...
		{
			name: "",
			args: []string{
				"--no-config-databases",
				filepath.FromSlash("./testdata/configs-extra-dbs/yarn.lock"),
			},
			exit: 0,
		},
		// ...but it does still use the built-in databases
		{
			name: "",
			args: []string{
				"--config", filepath.FromSlash("./testdata/configs-extra-dbs/.osv-detector.yaml"),
				"--no-config-databases",
				filepath.FromSlash("./testdata/locks-many/yarn.lock"),
			},
			exit: 0,
		},
		// when a global config is provided, any local configs should be ignored
		{
			name: "",
			args: []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-one/yarn.lock")},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-two")},
			exit: 0,
		},
		{
			name: "",
			args: []string{"--config", filepath.FromSlash("testdata/my-config.yml"), filepath.FromSlash("./testdata/configs-one/yarn.lock"), filepath.FromSlash("./testdata/locks-many")},
			exit: 0,
		},
		// when a local config is invalid, none of the lockfiles in that directory should
		// be checked (as the results could be different due to e.g. missing ignores)
		{
			name: "",
			args: []string{filepath.FromSlash("./testdata/configs-invalid"), filepath.FromSlash("./testdata/locks-one")},
			exit: 127,
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
			exit: 127,
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
			name: "",
			args: []string{"--ignore", "GHSA-1234", "--ignore", "GHSA-5678", filepath.FromSlash("./testdata/locks-one")},
			exit: 0,
		},
		{
			name: "",
			args: []string{
				"--ignore", "GHSA-whgm-jr23-g3j9",
				"--parse-as", "package-lock.json",
				filepath.FromSlash("./testdata/locks-insecure/my-package-lock.json"),
			},
			exit: 0,
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
			exit: 0,
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
			exit: 0,
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
			exit: 1,
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
			exit: 0,
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
			exit: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}

func setupConfigForUpdating(t *testing.T, path string, initial string) func() {
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

		snaps.MatchSnapshot(t, string(content))
	}
}

func TestRun_UpdatingConfigIgnores(t *testing.T) {
	t.Parallel()

	tests := []cliTestCase{
		// when there is no existing config, nothing should be updated
		{
			name: "",
			args: []string{"--update-config-ignores", filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json")},
			exit: 1,
		},
		// when given an explicit config, that should be updated
		{
			name: "",
			args: []string{
				"--update-config-ignores",
				"--config", "testdata/existing-config.yml",
				filepath.FromSlash("package-lock.json:./testdata/locks-insecure/my-package-lock.json"),
			},
			exit: 1,
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t, "testdata/existing-config.yml", "")
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
			exit: 0,
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-ignores.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
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
			exit: 1,
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-ignored-ignores.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
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
			exit: 1,
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/existing-config-with-many-lockfiles.yml",
					"ignore: [GHSA-whgm-jr23-g3j9]",
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
			exit: 1,
			around: func(t *testing.T) func() {
				t.Helper()

				cleanupConfig1 := setupConfigForUpdating(t,
					"testdata/locks-insecure-nested/.osv-detector.yml",
					"ignore: []",
				)

				cleanupConfig2 := setupConfigForUpdating(t,
					"testdata/locks-insecure-nested/nested/.osv-detector.yml",
					"ignore: []",
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
			exit: 1,
			around: func(t *testing.T) func() {
				t.Helper()

				return setupConfigForUpdating(t,
					"testdata/locks-insecure-many/.osv-detector.yml",
					"ignore: [GHSA-7p7h-4mm5-852v, GHSA-93q8-gq69-wqmw, GHSA-67hx-6x53-jw92]",
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
		snaps.Skip(t, "Skipping acceptance tests")
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

		tests = append(tests, cliTestCase{
			args: []string{parseAs + ":" + fp},
			exit: 1,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCli(t, tt)
		})
	}
}
