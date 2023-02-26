package lockfile

import (
	"errors"
	"fmt"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Represents when a package name could not be determined while parsing.
// Currently, parsers are expected to omit such packages from their results.
const unknownPkgName = "<unknown>"

func FindParser(pathToLockfile string, parseAs string) (PackageDetailsParser, string) {
	if parseAs == "" {
		parseAs = filepath.Base(pathToLockfile)
	}

	return parsers[parseAs], parseAs
}

//nolint:gochecknoglobals // this is an optimisation and read-only
var parsers = map[string]PackageDetailsParser{
	"buildscript-gradle.lockfile": ParseGradleLockFile,
	"Cargo.lock":                  ParseCargoLockFile,
	"composer.lock":               ParseComposerLockFile,
	"Gemfile.lock":                ParseGemfileLockFile,
	"go.mod":                      ParseGoLockFile,
	"gradle.lockfile":             ParseGradleLockFile,
	"mix.lock":                    ParseMixLockFile,
	"Pipfile.lock":                ParsePipenvLockFile,
	"package-lock.json":           ParseNpmLockFile,
	"packages.lock.json":          ParseNuGetLockFile,
	"pnpm-lock.yaml":              ParsePnpmLockFile,
	"poetry.lock":                 ParsePoetryLockFile,
	"pom.xml":                     ParseMavenLockFile,
	"pubspec.lock":                ParsePubspecLockFile,
	"requirements.txt":            ParseRequirementsTxtFile,
	"yarn.lock":                   ParseYarnLockFile,
}

//nolint:gochecknoglobals // this is an optimisation and read-only
var parsersWithReaders = map[string]PackageDetailsParserWithReader{
	"buildscript-gradle.lockfile": ParseGradleLock,
	"Cargo.lock":                  ParseCargoLock,
	"composer.lock":               ParseComposerLock,
	"Gemfile.lock":                ParseGemfileLock,
	"go.mod":                      ParseGoLock,
	"gradle.lockfile":             ParseGradleLock,
	"mix.lock":                    ParseMixLock,
	"Pipfile.lock":                ParsePipenvLock,
	"package-lock.json":           ParseNpmLock,
	"packages.lock.json":          ParseNuGetLock,
	"pnpm-lock.yaml":              ParsePnpmLock,
	"poetry.lock":                 ParsePoetryLock,
	"pom.xml":                     ParseMavenLock,
	"pubspec.lock":                ParsePubspecLock,
	"requirements.txt":            ParseRequirementsTxt,
	"yarn.lock":                   ParseYarnLock,
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
		if ecosystem == "" {
			continue
		}

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

func (l Lockfile) String() string {
	lines := make([]string, 0, len(l.Packages))

	for _, details := range l.Packages {
		ecosystem := details.Ecosystem

		if ecosystem == "" {
			ecosystem = "<unknown>"
		}

		ln := fmt.Sprintf("  %s: %s", ecosystem, details.Name)

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
		if parseAs != "" {
			return Lockfile{}, fmt.Errorf("%w, requested %s", ErrParserNotFound, parseAs)
		}

		return Lockfile{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	packages, err := parser(pathToLockfile)

	if err != nil && parseAs != "" {
		err = fmt.Errorf("(parsing as %s) %w", parsedAs, err)
	}

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

// ParseInImage attempts to extract a collection of package details from a lockfile
// that resides in a container image, using one of the native parsers.
//
// The parser is selected based on the name of the file, which can be overridden
// with the "parseAs" parameter.
func ParseInImage(pathToLockfile string, parseAs string, img image.Image) (Lockfile, error) {
	parsedAs := parseAs
	if parsedAs == "" {
		parsedAs = filepath.Base(pathToLockfile)
	}

	parser := parsersWithReaders[parsedAs]

	if parser == nil {
		if parseAs != "" {
			return Lockfile{}, fmt.Errorf("%w, requested %s", ErrParserNotFound, parseAs)
		}

		return Lockfile{}, fmt.Errorf("%w for %s", ErrParserNotFound, pathToLockfile)
	}

	r, err := img.OpenPathFromSquash(file.Path(pathToLockfile))

	if err != nil && parseAs != "" {
		return Lockfile{}, fmt.Errorf("(parsing as %s) %w", parsedAs, err)
	}

	packages, err := parser(r)

	if err != nil && parseAs != "" {
		err = fmt.Errorf("(parsing as %s) %w", parsedAs, err)
	}

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

func parseFile(pathToLockfile string, parserWithReader PackageDetailsParserWithReader) ([]PackageDetails, error) {
	r, err := os.Open(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	details, err := parserWithReader(r)

	if err != nil {
		err = fmt.Errorf("error while parsing %s: %w", pathToLockfile, err)
	}

	return details, err
}
