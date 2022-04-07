package reporter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"osv-detector/internal/reporter"
	"strings"
	"testing"
)

type TestResult struct {
	String               string `json:"value"`
	ErrorWhenMarshalling bool   `json:"-"`
}

func (r TestResult) ToString() string {
	return r.String
}

func (r TestResult) MarshalJSON() ([]byte, error) {
	type rawTestResult TestResult

	if r.ErrorWhenMarshalling {
		return nil, fmt.Errorf("oh noes, an error")
	}

	out, err := json.Marshal((rawTestResult)(r))

	if err != nil {
		return out, fmt.Errorf("%w", err)
	}

	return out, nil
}

func TestReporter_PrintError(t *testing.T) {
	t.Parallel()

	msg := "Hello world!"
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, false)

	r.PrintError(msg)

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
	expected := "{\"results\":[{\"value\":\"Hello world!\"}]}"
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

	if gotStdout := stdout.String(); gotStdout != expected {
		t.Errorf("Expected stdout to be \"%s\", but got \"%s\"", expected, gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != "" {
		t.Errorf("Expected stderr to be empty, but got \"%s\"", gotStderr)
	}
}

func TestReporter_PrintResult_OutputAsJSON_Error(t *testing.T) {
	t.Parallel()

	msg := "Hello sunshine!"
	expected := "an error occurred when printing results as JSON"
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	r := reporter.New(stdout, stderr, true)

	r.PrintResult(TestResult{String: msg, ErrorWhenMarshalling: true})

	if gotStdout := stdout.String(); gotStdout != "" {
		t.Errorf("Expected stdout to be empty, but got \"%s\"", gotStdout)
	}

	if gotStderr := stderr.String(); gotStderr != "" {
		t.Errorf("Expected stderr to be empty, but got \"%s\"", gotStderr)
	}

	r.PrintJSONResults()

	if gotStdout := stdout.String(); gotStdout != "" {
		t.Errorf("Expected stdout to be empty, but got \"%s\"", gotStdout)
	}

	if gotStderr := stderr.String(); !strings.Contains(gotStderr, expected) {
		t.Errorf("Expected stderr to contain \"%s\", but got \"%s\"", expected, gotStderr)
	}
}

func TestReporter_PrintText(t *testing.T) {
	t.Parallel()

	type args struct {
		outputAsJSON bool
		msg          string
	}
	tests := []struct {
		name         string
		args         args
		wantedStdout string
		wantedStderr string
	}{
		{
			name: "",
			args: args{
				outputAsJSON: false,
				msg:          "hello world",
			},
			wantedStdout: "hello world",
			wantedStderr: "",
		},
		{
			name: "",
			args: args{
				outputAsJSON: true,
				msg:          "hello world",
			},
			wantedStdout: "",
			wantedStderr: "hello world",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			r := reporter.New(stdout, stderr, tt.args.outputAsJSON)
			r.PrintText(tt.args.msg)

			if gotStdout := stdout.String(); gotStdout != tt.wantedStdout {
				t.Errorf("stdout got = %s, want %s", gotStdout, tt.wantedStdout)
			}

			if gotStderr := stderr.String(); gotStderr != tt.wantedStderr {
				t.Errorf("stderr got = %s, want %s", gotStderr, tt.wantedStderr)
			}
		})
	}
}
