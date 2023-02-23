package lockfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
)

const PipEcosystem Ecosystem = "PyPI"

// todo: expand this to support more things, e.g.
//
//	https://pip.pypa.io/en/stable/reference/requirements-file-format/#example
func parseLine(line string) PackageDetails {
	var constraint string
	name := line

	version := "0.0.0"

	if strings.Contains(line, "==") {
		constraint = "=="
	}

	if strings.Contains(line, ">=") {
		constraint = ">="
	}

	if strings.Contains(line, "~=") {
		constraint = "~="
	}

	if strings.Contains(line, "!=") {
		constraint = "!="
	}

	if constraint != "" {
		unprocessedName, unprocessedVersion, _ := strings.Cut(line, constraint)
		name = strings.TrimSpace(unprocessedName)

		if constraint != "!=" {
			version, _, _ = strings.Cut(strings.TrimSpace(unprocessedVersion), " ")
		}
	}

	return PackageDetails{
		Name:      normalizedRequirementName(name),
		Version:   version,
		Ecosystem: PipEcosystem,
		CompareAs: PipEcosystem,
	}
}

// normalizedName ensures that the package name is normalized per PEP-0503
// and then removing "added support" syntax if present.
//
// This is done to ensure we don't miss any advisories, as while the OSV
// specification says that the normalized name should be used for advisories,
// that's not the case currently in our databases, _and_ Pip itself supports
// non-normalized names in the requirements.txt, so we need to normalize
// on both sides to ensure we don't have false negatives.
//
// It's possible that this will cause some false positives, but that is better
// than false negatives, and can be dealt with when/if it actually happens.
func normalizedRequirementName(name string) string {
	// per https://www.python.org/dev/peps/pep-0503/#normalized-names
	name = cachedregexp.MustCompile(`[-_.]+`).ReplaceAllString(name, "-")
	name = strings.ToLower(name)
	name, _, _ = strings.Cut(name, "[")

	return name
}

func removeComments(line string) string {
	var re = cachedregexp.MustCompile(`(^|\s+)#.*$`)

	return strings.TrimSpace(re.ReplaceAllString(line, ""))
}

func isNotRequirementLine(line string) bool {
	return line == "" ||
		// flags are not supported
		strings.HasPrefix(line, "-") ||
		// file urls
		strings.HasPrefix(line, "https://") ||
		strings.HasPrefix(line, "http://") ||
		// file paths are not supported (relative or absolute)
		strings.HasPrefix(line, ".") ||
		strings.HasPrefix(line, "/")
}

func isLineContinuation(line string) bool {
	// checks that the line ends with an odd number of backslashes,
	// meaning the last one isn't escaped
	var re = cachedregexp.MustCompile(`([^\\]|^)(\\{2})*\\$`)

	return re.MatchString(line)
}

func ParseRequirementsTxtFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseRequirementsTxt)
}

func ParseRequirementsTxt(r io.Reader) ([]PackageDetails, error) {
	return parseRequirementsTxt(r /*, map[string]struct{}{}*/)
}

func parseRequirementsTxt(r io.Reader /*, requiredAlready map[string]struct{}*/) ([]PackageDetails, error) {
	packages := map[string]PackageDetails{}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		for isLineContinuation(line) {
			line = strings.TrimSuffix(line, "\\")

			if scanner.Scan() {
				line += scanner.Text()
			}
		}

		line = removeComments(line)

		// todo: figure out how to support this
		// if ar := strings.TrimPrefix(line, "-r "); ar != line {
		// 	ar = filepath.Join(filepath.Dir(pathToLockfile), ar)
		//
		// 	if _, ok := requiredAlready[ar]; ok {
		// 		continue
		// 	}
		//
		// 	requiredAlready[ar] = struct{}{}
		//
		// 	details, err := parseRequirementsTxt(ar, requiredAlready)
		//
		// 	if err != nil {
		// 		return []PackageDetails{}, fmt.Errorf("failed to include %s: %w", line, err)
		// 	}
		//
		// 	for _, detail := range details {
		// 		packages[detail.Name+"@"+detail.Version] = detail
		// 	}
		//
		// 	continue
		// }

		if isNotRequirementLine(line) {
			continue
		}

		detail := parseLine(line)
		packages[detail.Name+"@"+detail.Version] = detail
	}

	if err := scanner.Err(); err != nil {
		return []PackageDetails{}, fmt.Errorf("error while scanning: %w", err)
	}

	return pkgDetailsMapToSlice(packages), nil
}
