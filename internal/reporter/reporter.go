package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"osv-detector/pkg/database"
	"osv-detector/pkg/lockfile"
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

// PrintError writes the given message to stderr, regardless of if the reporter
// is outputting as JSON or not
func (r *Reporter) PrintError(msg string) {
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
	String() string
}

func (r *Reporter) PrintResult(result Result) {
	if r.outputAsJSON {
		r.results = append(r.results, result)

		return
	}

	fmt.Fprint(r.stdout, result.String())
}

// PrintJSONResults prints any results that this reporter has collected to
// stdout as JSON.
func (r *Reporter) PrintJSONResults() {
	out, err := json.Marshal(struct {
		Results interface{} `json:"results"`
	}{Results: r.results})

	if err != nil {
		r.PrintError(fmt.Sprintf("an error occurred when printing results as JSON: %v", err))

		return
	}

	fmt.Fprint(r.stdout, string(out))
}

func (r *Reporter) PrintDatabaseLoadErr(err error) {
	msg := err.Error()

	if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		msg = color.RedString("no local version of the database was found, and --offline flag was set")
	}

	r.PrintError(fmt.Sprintf(" %s\n", color.RedString("failed: %s", msg)))
}

func (r *Reporter) PrintKnownEcosystems() {
	ecosystems := lockfile.KnownEcosystems()

	r.PrintText("The detector supports parsing for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		r.PrintText(fmt.Sprintf("  %s\n", ecosystem))
	}
}
