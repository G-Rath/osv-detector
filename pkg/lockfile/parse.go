package lockfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

func FindParser(pathToLockfile string, parseAs string) (PackageDetailsParser, string) {
	if parseAs == "" {
		parseAs = path.Base(pathToLockfile)
	}

	return parsers[parseAs], parseAs
}

// nolint:gochecknoglobals // this is an optimisation and read-only
var parsers = map[string]PackageDetailsParser{
	"cargo.lock":        ParseCargoLock,
	"composer.lock":     ParseComposerLock,
	"Gemfile.lock":      ParseGemfileLock,
	"go.mod":            ParseGoLock,
	"package-lock.json": ParseNpmLock,
	"pnpm-lock.yaml":    ParsePnpmLock,
	"pom.xml":           ParseMavenLock,
	"requirements.txt":  ParseRequirementsTxt,
	"yarn.lock":         ParseYarnLock,
}

func ListParsers() []string {
	ps := make([]string, 0, len(parsers))

	for s := range parsers {
		ps = append(ps, s)
	}

	sort.Slice(ps, func(i, j int) bool {
		return strings.ToLower(ps[i]) < strings.ToLower(ps[j])
	})

	return ps
}

var ErrParserNotFound = errors.New("could not determine parser")

type Packages []PackageDetails

func toSliceOfEcosystems(ecosystemsMap map[Ecosystem]struct{}) []Ecosystem {
	ecosystems := make([]Ecosystem, 0, len(ecosystemsMap))

	for ecosystem := range ecosystemsMap {
		ecosystems = append(ecosystems, ecosystem)
	}

	return ecosystems
}

func (ps Packages) Ecosystems() []Ecosystem {
	ecosystems := make(map[Ecosystem]struct{})

	for _, pkg := range ps {
		ecosystems[pkg.Ecosystem] = struct{}{}
	}

	slicedEcosystems := toSliceOfEcosystems(ecosystems)

	sort.Slice(slicedEcosystems, func(i, j int) bool {
		return slicedEcosystems[i] < slicedEcosystems[j]
	})

	return slicedEcosystems
}

type Lockfile struct {
	FilePath string   `json:"filePath"`
	ParsedAs string   `json:"parsedAs"`
	Packages Packages `json:"packages"`
}

func (l Lockfile) ToString() string {
	lines := make([]string, 0, len(l.Packages))

	for _, details := range l.Packages {
		ln := fmt.Sprintf("  %s: %s", details.Ecosystem, details.Name)

		if details.Version != "" {
			ln += "@" + details.Version
		}

		if details.Commit != "" {
			ln += " (" + details.Commit + ")"
		}

		lines = append(lines, ln)
	}

	return strings.Join(lines, "\n")
}

// Parse attempts to extract a collection of package details from a lockfile,
// using one of the native parsers.
//
// The parser is selected based on the name of the file, which can be overridden
// with the "parseAs" parameter.
func Parse(pathToLockfile string, parseAs string) (Lockfile, error) {
	parser, parsedAs := FindParser(pathToLockfile, parseAs)

	if parser == nil {
		return Lockfile{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	packages, err := parser(pathToLockfile)

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == packages[j].Name {
			return packages[i].Version < packages[j].Version
		}

		return packages[i].Name < packages[j].Name
	})

	return Lockfile{
		FilePath: pathToLockfile,
		ParsedAs: parsedAs,
		Packages: packages,
	}, err
}

var errCSVRecordNotEnoughFields = errors.New("not enough fields (missing at least ecosystem and package name)")
var errCSVRecordMissingEcosystemField = errors.New("field 1 is empty (must be the name of an ecosystem)")
var errCSVRecordMissingPackageField = errors.New("field 2 is empty (must be the name of a package)")

func fromCSVRecord(lines []string) (PackageDetails, error) {
	if len(lines) < 2 {
		return PackageDetails{}, errCSVRecordNotEnoughFields
	}

	if lines[0] == "" {
		return PackageDetails{}, errCSVRecordMissingEcosystemField
	}

	if lines[1] == "" {
		return PackageDetails{}, errCSVRecordMissingPackageField
	}

	return PackageDetails{
		Name:      lines[1],
		Version:   lines[2],
		Ecosystem: Ecosystem(lines[0]),
	}, nil
}

func fromCSV(reader io.Reader) ([]PackageDetails, error) {
	var packages []PackageDetails

	i := 0
	r := csv.NewReader(reader)

	for {
		i++
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return packages, fmt.Errorf("%w", err)
		}

		details, err := fromCSVRecord(record)
		if err != nil {
			return packages, fmt.Errorf("row %d: %w", i, err)
		}

		packages = append(packages, details)
	}

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == packages[j].Name {
			return packages[i].Version < packages[j].Version
		}

		return packages[i].Name < packages[j].Name
	})

	return packages, nil
}

func FromCSVRows(filePath string, parseAs string, rows []string) (Lockfile, error) {
	packages, err := fromCSV(strings.NewReader(strings.Join(rows, "\n")))

	return Lockfile{
		FilePath: filePath,
		ParsedAs: parseAs,
		Packages: packages,
	}, err
}

func FromCSVFile(pathToCSV string, parseAs string) (Lockfile, error) {
	reader, err := os.Open(pathToCSV)

	if err != nil {
		return Lockfile{}, fmt.Errorf("could not read %s: %w", pathToCSV, err)
	}

	packages, err := fromCSV(reader)

	if err != nil {
		err = fmt.Errorf("%s: %w", pathToCSV, err)
	}

	return Lockfile{
		FilePath: pathToCSV,
		ParsedAs: parseAs,
		Packages: packages,
	}, err
}
