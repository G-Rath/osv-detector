package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/lockfile"
)

type Reporter struct {
	stdout       io.Writer
	stderr       io.Writer
	outputAsJSON bool
	results      []Result

	hasErrored bool
}

func New(stdout io.Writer, stderr io.Writer, outputAsJSON bool) *Reporter {
	return &Reporter{
		stdout:       stdout,
		stderr:       stderr,
		outputAsJSON: outputAsJSON,
		results:      make([]Result, 0),
	}
}

func (r *Reporter) HasErrored() bool {
	return r.hasErrored
}

// PrintErrorf writes the given message to stderr, regardless of if the reporter
// is outputting as JSON or not
func (r *Reporter) PrintErrorf(msg string, a ...any) {
	r.hasErrored = true

	// todo: this is a hack to make the lockfile/extractor error output more like the original
	//  there's no real reason to be doing it other than that there isn't a reason not to...
	fmt.Fprint(r.stderr, strings.Replace(fmt.Sprintf(msg, a...), " could not extract:", "", 1))
}

// PrintTextf writes the given message to stdout, _unless_ the reporter is set
// to output as JSON, in which case it writes the message to stderr.
//
// This should be used for content that should always be outputted, but that
// should not be captured when piping if outputting JSON.
func (r *Reporter) PrintTextf(msg string, a ...any) {
	target := r.stdout

	if r.outputAsJSON {
		target = r.stderr
	}

	fmt.Fprintf(target, msg, a...)
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
		Results any `json:"results"`
	}{Results: r.results})

	if err != nil {
		r.PrintErrorf("an error occurred when printing results as JSON: %v", err)

		return
	}

	fmt.Fprint(r.stdout, string(out))
}

func (r *Reporter) PrintDatabaseLoadErr(err error) {
	msg := err.Error()

	if errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		msg = color.RedString("no local version of the database was found, and --offline flag was set")
	}

	r.PrintErrorf(" %s\n", color.RedString("failed: %s", msg))
}

func (r *Reporter) PrintKnownEcosystems() {
	ecosystems := lockfile.KnownEcosystems()

	r.PrintTextf("The detector supports parsing for the following ecosystems:\n")

	for _, ecosystem := range ecosystems {
		r.PrintTextf("  %s\n", ecosystem)
	}
}
