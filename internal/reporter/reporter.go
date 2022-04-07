package reporter

import (
	"encoding/json"
	"fmt"
	"io"
)

type Reporter struct {
	stdout       io.Writer
	stderr       io.Writer
	outputAsJSON bool
	results      []Result
}

func New(stdout io.Writer, stderr io.Writer, outputAsJSON bool) *Reporter {
	return &Reporter{
		stdout:       stdout,
		stderr:       stderr,
		outputAsJSON: outputAsJSON,
		results:      make([]Result, 0),
	}
}

// PrintExtra writes the given message to stderr
func (r *Reporter) PrintExtra(msg string) {
	fmt.Fprint(r.stderr, msg)
}

// PrintText writes the given message to stdout, _unless_ the reporter is set
// to output as JSON, in which case it writes the message to stderr.
//
// This should be used for content that should always be outputted, but that
// should not be captured when piping if outputting JSON.
func (r *Reporter) PrintText(msg string) {
	target := r.stdout

	if r.outputAsJSON {
		target = r.stderr
	}

	fmt.Fprint(target, msg)
}

type Result interface {
	ToString() string
}

func (r *Reporter) PrintResult(result Result) {
	if r.outputAsJSON {
		r.results = append(r.results, result)

		return
	}

	fmt.Fprint(r.stdout, result.ToString())
}

// PrintJSONResults prints any results that this reporter has collected to
// stdout as JSON.
func (r *Reporter) PrintJSONResults() {
	out, err := json.Marshal(struct {
		Results interface{} `json:"results"`
	}{Results: r.results})

	if err != nil {
		r.PrintExtra(fmt.Sprintf("an error occurred when printing results as JSON: %v", err))

		return
	}

	fmt.Fprint(r.stdout, string(out))
}
