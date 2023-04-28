package lockfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isNotGradleDependencyLine(line string) bool {
	return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "empty=")
}

func parseGradleLine(line string) (PackageDetails, error) {
	parts := strings.SplitN(line, ":", 3)
	if len(parts) < 3 {
		return PackageDetails{}, fmt.Errorf("invalid line in gradle lockfile: %s", line) //nolint:goerr113
	}

	group, artifact, version := parts[0], parts[1], parts[2]
	version, _, _ = strings.Cut(version, "=")

	return PackageDetails{
		Name:      fmt.Sprintf("%s:%s", group, artifact),
		Version:   version,
		Ecosystem: MavenEcosystem,
		CompareAs: MavenEcosystem,
	}, nil
}

func ParseGradleLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseGradleLock)
}

func ParseGradleLock(f ParsableFile) ([]PackageDetails, error) {
	var packages []PackageDetails

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lockLine := strings.TrimSpace(scanner.Text())

		if isNotGradleDependencyLine(lockLine) {
			continue
		}

		pkg, err := parseGradleLine(lockLine)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %s\n", err.Error())

			continue
		}

		packages = append(packages, pkg)
	}

	if err := scanner.Err(); err != nil {
		return []PackageDetails{}, fmt.Errorf("failed to read: %w", err)
	}

	return packages, nil
}
