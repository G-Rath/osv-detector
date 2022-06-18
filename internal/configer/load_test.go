package configer_test

import (
	"bytes"
	"fmt"
	"osv-detector/internal/configer"
	"osv-detector/internal/reporter"
	"reflect"
	"strings"
	"testing"
)

// nolint:unparam
func newReporter(t *testing.T) (*reporter.Reporter, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, false)

	return r, stdout, stderr
}

func TestFind_NoConfig(t *testing.T) {
	t.Parallel()

	r, _, _ := newReporter(t)

	config, err := configer.Find(r, "fixtures/no-config")

	if err != nil {
		t.Errorf("Find() error = %v, expected nothing", err)
	}

	if config.FilePath != "" {
		t.Errorf("Find() config.FilePath = %s, expected empty", config.FilePath)
	}

	if l := len(config.Ignore); l > 0 {
		t.Errorf("Find() config.Ignore = %d, expected empty", l)
	}
}

func TestFind_ExtYml(t *testing.T) {
	t.Parallel()

	r, _, _ := newReporter(t)

	expectedIgnores := []string{"GHSA-1", "GHSA-2", "GHSA-3"}
	expectedFilePath := "fixtures/ext-yml/.osv-detector.yml"

	config, err := configer.Find(r, "fixtures/ext-yml")

	if err != nil {
		t.Errorf("Find() error = %v, expected nothing", err)
	}

	if config.FilePath != expectedFilePath {
		t.Errorf("Find() config.FilePath = %s, expected %s", config.FilePath, expectedFilePath)
	}

	if !reflect.DeepEqual(config.Ignore, expectedIgnores) {
		t.Errorf("Find() config.Ignore = %v, expected empty", expectedIgnores)
	}
}

func TestFind_ExtYml_Invalid(t *testing.T) {
	t.Parallel()

	r, _, _ := newReporter(t)

	expectedFilePath := "fixtures/ext-yml-invalid/.osv-detector.yml"

	config, err := configer.Find(r, "fixtures/ext-yml-invalid")

	if err == nil {
		t.Errorf("Find() did not error, which was unexpected")
	}

	if config.FilePath != expectedFilePath {
		t.Errorf("Find() config.FilePath = %s, expected %s", config.FilePath, expectedFilePath)
	}

	if l := len(config.Ignore); l > 0 {
		t.Errorf("Find() config.Ignore = %d, expected empty", l)
	}
}

func TestFind_ExtYaml(t *testing.T) {
	t.Parallel()

	r, _, _ := newReporter(t)

	expectedIgnores := []string{"GHSA-4", "GHSA-5", "GHSA-6"}
	expectedFilePath := "fixtures/ext-yaml/.osv-detector.yaml"

	config, err := configer.Find(r, "fixtures/ext-yaml")

	if err != nil {
		t.Errorf("Find() error = %v, expected nothing", err)
	}

	if config.FilePath != expectedFilePath {
		t.Errorf("Find() config.FilePath = %s, expected %s", config.FilePath, expectedFilePath)
	}

	if !reflect.DeepEqual(config.Ignore, expectedIgnores) {
		t.Errorf("Find() config.Ignore = %v, expected empty", expectedIgnores)
	}
}

func TestFind_ExtYaml_Invalid(t *testing.T) {
	t.Parallel()

	r, _, _ := newReporter(t)

	expectedFilePath := "fixtures/ext-yaml-invalid/.osv-detector.yaml"

	config, err := configer.Find(r, "fixtures/ext-yaml-invalid")

	if err == nil {
		t.Errorf("Find() did not error, which was unexpected")
	}

	if config.FilePath != expectedFilePath {
		t.Errorf("Find() config.FilePath = %s, expected %s", config.FilePath, expectedFilePath)
	}

	if l := len(config.Ignore); l > 0 {
		t.Errorf("Find() config.Ignore = %d, expected empty", l)
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()

	r, _, stderr := newReporter(t)

	fixturePath := "fixtures/extra-databases/.osv-detector.yml"

	config, err := configer.Load(r, fixturePath)

	if err != nil {
		t.Fatalf("Load() unexpected errored \"%v\"", err)
	}

	if len(config.Databases) != 7 {
		t.Fatalf("Load() returned a config that has %d databases instead of 7", len(config.Databases))
	}

	gotStderr := stderr.String()

	expectedErrors := []string{
		"unsupported database source type file",
		"invalid URI",
	}

	isMissingStderrMessage := false
	for _, expectedError := range expectedErrors {
		if !strings.Contains(gotStderr, expectedError) {
			t.Errorf("Expected stderr to contain \"%s\" but it did not", expectedError)
			isMissingStderrMessage = true
		}
	}

	if isMissingStderrMessage {
		// nolint:forbidigo
		fmt.Println(gotStderr)
	}

	if config.FilePath != fixturePath {
		t.Errorf("Find() config.FilePath = %s, expected %s", config.FilePath, fixturePath)
	}

	if l := len(config.Ignore); l > 0 {
		t.Errorf("Find() config.Ignore = %d, expected empty", l)
	}

	expectedTypes := []string{
		"zip",
		"dir",
		"dir",
		"api",
		"zip",
		"zip",
		"zip",
	}

	expectedNames := []string{
		"zip#https://github.com/github/advisory-database/archive/refs/heads/main.zip",
		"dir#file:/relative/path/to/dir",
		"dir#file:////root/path/to/dir",
		"api#https://api-staging.osv.dev/v1",
		"GitHub Advisory Database",
		"zip#https://my-site.com/osvs/all",
		"zip#https://github.com/github/advisory-database/archive/refs/heads/main.zip#advisory-database-main/advisories/unreviewed",
	}

	for i, db := range config.Databases {
		if db.Type != expectedTypes[i] {
			t.Errorf("Load() expected db %d to be type %s but got %s", i, expectedTypes[i], db.Type)
		}

		if db.Name != expectedNames[i] {
			t.Errorf("Load() expected db %d to be named %s but got %s", i, expectedNames[i], db.Name)
		}
	}

	lastDatabase := config.Databases[len(config.Databases)-1]

	if lastDatabase.WorkingDirectory != "advisory-database-main/advisories/unreviewed" {
		t.Errorf("Load() config had incorrect WorkingDirectory (\"%s\")", lastDatabase.WorkingDirectory)
	}
}
