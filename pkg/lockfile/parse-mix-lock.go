package lockfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
)

const MixEcosystem Ecosystem = "Hex"

func ParseMixLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseMixLock)
}

func ParseMixLock(r io.Reader) ([]PackageDetails, error) {
	re := cachedregexp.MustCompile(`^ +"(\w+)": \{.+,$`)

	scanner := bufio.NewScanner(r)

	var packages []PackageDetails

	for scanner.Scan() {
		line := scanner.Text()

		match := re.FindStringSubmatch(line)

		if match == nil {
			continue
		}

		// we only care about the third and fourth "rows" which are both strings,
		// so we can safely split the line as if it's a set of comma-separated fields
		// even though that'll actually poorly represent nested arrays & objects
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ','
		})

		if len(fields) < 4 {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"Found less than four fields when parsing a line that looks like a dependency in a mix.lock - please report this!\n",
			)

			continue
		}

		name := match[1]
		version := strings.TrimSpace(fields[2])
		commit := strings.TrimSpace(fields[3])

		version = strings.TrimSuffix(strings.TrimPrefix(version, `"`), `"`)
		commit = strings.TrimSuffix(strings.TrimPrefix(commit, `"`), `"`)

		if strings.HasSuffix(fields[0], ":git") {
			commit = version
			version = ""
		}

		packages = append(packages, PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: MixEcosystem,
			CompareAs: MixEcosystem,
			Commit:    commit,
		})
	}

	if err := scanner.Err(); err != nil {
		return []PackageDetails{}, fmt.Errorf("error while scanning: %w", err)
	}

	return packages, nil
}
