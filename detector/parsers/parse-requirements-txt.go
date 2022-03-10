package parsers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const PipEcosystem Ecosystem = "PyPI"

// todo: expand this to support more things, e.g.
//   https://pip.pypa.io/en/stable/reference/requirements-file-format/#example
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
		splitted := strings.Split(line, constraint)

		name = strings.TrimSpace(splitted[0])

		if constraint != "!=" {
			version = strings.TrimSpace(splitted[1])
		}
	}

	return PackageDetails{
		Name:      cleanupRequirementName(name),
		Version:   version,
		Ecosystem: PipEcosystem,
	}
}

func cleanupRequirementName(name string) string {
	return strings.Split(name, "[")[0]
}

func removeComments(line string) string {
	var re = regexp.MustCompile(`(^|\s+)#.*$`)

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

func ParseRequirementsTxt(pathToLockfile string) ([]PackageDetails, error) {
	var packages []PackageDetails

	file, err := os.Open(pathToLockfile)
	if err != nil {
		return packages, fmt.Errorf("could not open %s: %w", pathToLockfile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := removeComments(scanner.Text())

		if isNotRequirementLine(line) {
			continue
		}

		packages = append(packages, parseLine(line))
	}

	if err := scanner.Err(); err != nil {
		return packages, fmt.Errorf("error while scanning %s: %w", pathToLockfile, err)
	}

	return packages, nil
}
