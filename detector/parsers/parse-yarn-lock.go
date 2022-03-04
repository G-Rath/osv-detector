package parsers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const YarnEcosystem = NpmEcosystem

func shouldSkipYarnLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func groupPackageLines(scanner *bufio.Scanner) [][]string {
	var groups [][]string
	var group []string

	for scanner.Scan() {
		line := scanner.Text()

		if shouldSkipYarnLine(line) {
			continue
		}

		// represents the start of a new dependency
		if !strings.HasPrefix(line, " ") {
			if len(group) > 0 {
				groups = append(groups, group)
			}
			group = make([]string, 0)
		}

		group = append(group, line)
	}

	if len(group) > 0 {
		groups = append(groups, group)
	}

	return groups
}

func extractYarnPackageName(str string) string {
	str = strings.TrimPrefix(str, "\"")

	isScoped := strings.HasPrefix(str, "@")

	if isScoped {
		str = strings.TrimPrefix(str, "@")
	}

	name := strings.SplitN(str, "@", 2)[0]

	if isScoped {
		name = "@" + name
	}

	return name
}

func determineYarnPackageVersion(group []string) string {
	re := regexp.MustCompile(`^ {2}version:? "?([\d.]+)"?$`)

	for _, s := range group {
		matched := re.FindStringSubmatch(s)

		if matched != nil {
			return matched[1]
		}
	}

	// todo: decide what to do here - maybe panic...?
	return ""
}

func parsePackageGroup(group []string) PackageDetails {
	return PackageDetails{
		Name:      extractYarnPackageName(group[0]),
		Version:   determineYarnPackageVersion(group),
		Ecosystem: YarnEcosystem,
	}
}

func ParseYarnLock(pathToLockfile string) ([]PackageDetails, error) {
	var packages []PackageDetails

	file, err := os.Open(pathToLockfile)
	if err != nil {
		return packages, fmt.Errorf("could not open %s: %w", pathToLockfile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	packageGroups := groupPackageLines(scanner)

	if err := scanner.Err(); err != nil {
		return packages, fmt.Errorf("error while scanning %s: %w", pathToLockfile, err)
	}

	for _, group := range packageGroups {
		if group[0] == "__metadata:" {
			continue
		}

		packages = append(packages, parsePackageGroup(group))
	}

	return packages, nil
}
