package reporter_test

import (
	"bytes"
	"osv-detector/internal/reporter"
	"testing"
)

type TestResult struct {
	String string `json:"value"`
}

func (r TestResult) ToString() string {
	return r.String
}

func TestReporter_PrintExtra(t *testing.T) {
	t.Parallel()

	msg := "Hello world!"
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, false)

	r.PrintExtra(msg)

	if gotStdout := stdout.String(); gotStdout != "" {
		t.Errorf("Expected stdout to be empty, but got \"%s\"", gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != msg {
		t.Errorf("Expected stderr to have \"%s\", but got \"%s\"", msg, gotStderr)
	}
}

func TestReporter_PrintResult(t *testing.T) {
	t.Parallel()

	msg := "Hello sunshine!"
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, false)

	r.PrintResult(TestResult{String: msg})

	if gotStdout := stdout.String(); gotStdout != msg {
		t.Errorf("Expected stdout to have \"%s\", but got \"%s\"", msg, gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != "" {
		t.Errorf("Expected stderr to be empty, but got \"%s\"", gotStderr)
	}
}

func TestReporter_PrintResult_OutputAsJSON(t *testing.T) {
	t.Parallel()

	msg := "Hello world!"
	json := "{\"results\":[{\"value\":\"Hello world!\"}]}"
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, true)

	r.PrintResult(TestResult{String: msg})

	if gotStdout := stdout.String(); gotStdout != "" {
		t.Errorf("Expected stdout to be empty, but got \"%s\"", gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != "" {
		t.Errorf("Expected stderr to be empty, but got \"%s\"", gotStderr)
	}

	r.PrintJSONResults()

	if gotStdout := stdout.String(); gotStdout != json {
		t.Errorf("Expected stdout to be \"%s\", but got \"%s\"", json, gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != "" {
		t.Errorf("Expected stderr to be empty, but got \"%s\"", gotStderr)
	}
}
