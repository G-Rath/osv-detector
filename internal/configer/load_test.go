package configer_test

import (
	"osv-detector/internal/configer"
	"reflect"
	"testing"
)

func TestFind_NoConfig(t *testing.T) {
	t.Parallel()

	config, err := configer.Find("fixtures/no-config")

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

	expectedIgnores := []string{"GHSA-1", "GHSA-2", "GHSA-3"}
	expectedFilePath := "fixtures/ext-yml/.osv-detector.yml"

	config, err := configer.Find("fixtures/ext-yml")

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

	expectedFilePath := "fixtures/ext-yml-invalid/.osv-detector.yml"

	config, err := configer.Find("fixtures/ext-yml-invalid")

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

	expectedIgnores := []string{"GHSA-4", "GHSA-5", "GHSA-6"}
	expectedFilePath := "fixtures/ext-yaml/.osv-detector.yaml"

	config, err := configer.Find("fixtures/ext-yaml")

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

	expectedFilePath := "fixtures/ext-yaml-invalid/.osv-detector.yaml"

	config, err := configer.Find("fixtures/ext-yaml-invalid")

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
