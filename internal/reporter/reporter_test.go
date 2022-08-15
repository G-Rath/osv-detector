package reporter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/g-rath/osv-detector/internal/reporter"
	"github.com/g-rath/osv-detector/pkg/database"
	"strings"
	"testing"
)

type TestResult struct {
	Value                string `json:"value"`
	ErrorWhenMarshalling bool   `json:"-"`
}

func (r TestResult) String() string {
	return r.Value
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

	r.PrintResult(TestResult{Value: msg})

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

	r.PrintResult(TestResult{Value: msg})

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

	r.PrintResult(TestResult{Value: msg, ErrorWhenMarshalling: true})

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

func TestReporter_PrintDatabaseLoadErr(t *testing.T) {
	t.Parallel()

	type args struct {
		outputAsJSON bool
		err          error
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
				err:          fmt.Errorf("oh noes"),
			},
			wantedStdout: "",
			wantedStderr: " failed: oh noes\n",
		},
		{
			name: "",
			args: args{
				outputAsJSON: true,
				err:          fmt.Errorf("oh noes"),
			},
			wantedStdout: "",
			wantedStderr: " failed: oh noes\n",
		},
		{
			name: "",
			args: args{
				outputAsJSON: false,
				err:          database.ErrOfflineDatabaseNotFound,
			},
			wantedStdout: "",
			wantedStderr: " failed: no local version of the database was found, and --offline flag was set\n",
		},
		{
			name: "",
			args: args{
				outputAsJSON: true,
				err:          database.ErrOfflineDatabaseNotFound,
			},
			wantedStdout: "",
			wantedStderr: " failed: no local version of the database was found, and --offline flag was set\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			r := reporter.New(stdout, stderr, tt.args.outputAsJSON)
			r.PrintDatabaseLoadErr(tt.args.err)

			if gotStdout := stdout.String(); gotStdout != tt.wantedStdout {
				t.Errorf("stdout got = \"%s\", want \"%s\"", gotStdout, tt.wantedStdout)
			}

			if gotStderr := stderr.String(); gotStderr != tt.wantedStderr {
				t.Errorf("stderr got = \"%s\", want \"%s\"", gotStderr, tt.wantedStderr)
			}
		})
	}
}

func TestReporter_PrintKnownEcosystems(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"The detector supports parsing for the following ecosystems:",
		"  npm",
		"  crates.io",
		"  RubyGems",
		"  Packagist",
		"  Go",
		"  Hex",
		"  Maven",
		"  PyPI",
		"",
	}, "\n")

	type args struct {
		outputAsJSON bool
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
			},
			wantedStdout: expected,
			wantedStderr: "",
		},
		{
			name: "",
			args: args{
				outputAsJSON: true,
			},
			wantedStdout: "",
			wantedStderr: expected,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			r := reporter.New(stdout, stderr, tt.args.outputAsJSON)
			r.PrintKnownEcosystems()

			if gotStdout := stdout.String(); gotStdout != tt.wantedStdout {
				t.Errorf("stdout got = \"%s\", want \"%s\"", gotStdout, tt.wantedStdout)
			}

			if gotStderr := stderr.String(); gotStderr != tt.wantedStderr {
				t.Errorf("stderr got = \"%s\", want \"%s\"", gotStderr, tt.wantedStderr)
			}
		})
	}
}
