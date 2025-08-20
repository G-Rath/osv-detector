package configer_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
	"github.com/g-rath/osv-detector/internal/configer"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
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

func createTestConfigFile(t *testing.T, content string) (string, func()) {
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "osv-detector-test-*")
	if err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	_, err = f.WriteString(dedent(t, content))

	if err != nil {
		t.Fatalf("could not write test config file file: %v", err)
	}

	return f.Name(), func() {
		_ = os.Remove(f.Name())
	}
}

func readConfigFixture(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("could not read config fixture file: %v", err)
	}

	return string(content)
}

func TestUpdateWithIgnores_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	err := configer.UpdateWithIgnores("testdata/does/not/exist", []string{})

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	if !strings.Contains(err.Error(), "could not read") {
		t.Errorf("Expected to get \"%s\" error, but got \"%v\"", "could not read", err)
	}
}

func expectAreEqual(t *testing.T, actual, expect string) {
	t.Helper()

	actual = dedent(t, actual)
	expect = dedent(t, expect)

	if !cmp.Equal(actual, expect) {
		if os.Getenv("TEST_NO_DIFF") == "true" {
			t.Errorf("\nactual does not match expected:\n got:\n%s\n\n want:\n%s", actual, expect)
		} else {
			t.Errorf("\nactual does not match expected:\n%s", cmp.Diff(expect, actual))
		}
	}
}

func TestUpdateWithIgnores(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ignores []string
		initial string
		updated string
	}{
		{
			name:    "completely empty file",
			ignores: []string{},
			initial: "",
			updated: "ignore: []",
		},
		{
			name:    "file with no ignores",
			ignores: []string{},
			initial: "ignore:",
			updated: "ignore: []",
		},
		{
			name:    "file with no ignores (adding some ignores)",
			ignores: []string{"OSV-1", "OSV-2", "OSV-3"},
			initial: "ignore:",
			updated: `
				ignore:
					- OSV-1
					- OSV-2
					- OSV-3
			`,
		},
		{
			name:    "ignores are not sorted or deduplicated",
			ignores: []string{"OSV-2", "OSV-5", "OSV-1", "OSV-2", "OSV-4", "OSV-3"},
			initial: "ignore:",
			updated: `
				ignore:
					- OSV-2
					- OSV-5
					- OSV-1
					- OSV-2
					- OSV-4
					- OSV-3
			`,
		},
		{
			name:    "comments and existing ignores are not preserved",
			ignores: []string{"OSV-1", "OSV-2"},
			initial: readConfigFixture(t, "testdata/ext-yml/.osv-detector.yml"),
			updated: `
				ignore:
					- OSV-1
					- OSV-2
			`,
		},
		{
			name:    "other config options are preserved",
			ignores: []string{"OSV-4", "OSV-5"},
			initial: readConfigFixture(t, "testdata/extra-databases/.osv-detector.yml"),
			updated: `
				ignore:
					- OSV-4
					- OSV-5
				extra-databases:
					- url: https://github.com/github/advisory-database/archive/refs/heads/main.zip
					- url: file:/relative/path/to/dir
					- url: file:////root/path/to/dir
					- url: https://api-staging.osv.dev/v1
					- url: https://github.com/github/advisory-database/archive/refs/heads/main.zip
						name: GitHub Advisory Database
					- url: https://my-site.com/osvs/all
						type: zip
					- url: https://github.com/github/advisory-database/archive/refs/heads/main.zip
						working-directory: advisory-database-main/advisories/unreviewed
					- url: https://my-site.com/osvs/all
						type: file
					- url: www.github.com/github/advisory-database/archive/refs/heads/main.zip
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			configFilePath, cleanupConfigFile := createTestConfigFile(t, tt.initial)
			defer cleanupConfigFile()

			err := configer.UpdateWithIgnores(configFilePath, tt.ignores)

			if err != nil {
				t.Fatalf("could not update config file: %v", err)
			}

			content, err := os.ReadFile(configFilePath)

			if err != nil {
				t.Fatalf("could not read test config file: %v", err)
			}

			expectAreEqual(t, string(content), tt.updated)

			if err = yaml.Unmarshal(content, &struct{}{}); err != nil {
				t.Fatalf("could not parse updated config as YAML: %v", err)
			}
		})
	}
}
