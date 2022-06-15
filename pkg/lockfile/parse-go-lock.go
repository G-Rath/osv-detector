package lockfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

const GoEcosystem Ecosystem = "Go"

func deduplicatePackages(packages map[string]PackageDetails) map[string]PackageDetails {
	details := map[string]PackageDetails{}

	for _, detail := range packages {
		details[detail.Name+"@"+detail.Version] = detail
	}

	return details
}

func parseGoModVersion(version string) (string, string, string) {
	re := regexp.MustCompile(`^v([\d.]+(?:-[\w.]+)?)[-.](\d{14})-(\w{12})$`)

	matched := re.FindStringSubmatch(version)

	if matched == nil {
		return version, "", ""
	}

	return matched[1], matched[2], matched[3]
}

func ParseGoLock(pathToLockfile string) ([]PackageDetails, error) {
	lockfileContents, err := os.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	parsedLockfile, err := modfile.Parse(pathToLockfile, lockfileContents, nil)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := map[string]PackageDetails{}

	for _, require := range parsedLockfile.Require {
		version, _, commit := parseGoModVersion(require.Mod.Version)

		packages[require.Mod.Path+"@"+require.Mod.Version] = PackageDetails{
			Name:      require.Mod.Path,
			Version:   strings.TrimPrefix(version, "v"),
			Ecosystem: GoEcosystem,
			CompareAs: GoEcosystem,
			Commit:    commit,
		}
	}

	for _, replace := range parsedLockfile.Replace {
		var replacements []string

		if replace.Old.Version == "" {
			// If the left version is omitted, all versions of the module are replaced.
			for k, pkg := range packages {
				if pkg.Name == replace.Old.Path {
					replacements = append(replacements, k)
				}
			}
		} else {
			// If a version is present on the left side of the arrow (=>),
			// only that specific version of the module is replaced
			s := replace.Old.Path + "@" + replace.Old.Version

			// A `replace` directive has no effect if the module version on the left side is not required.
			if _, ok := packages[s]; ok {
				replacements = []string{s}
			}
		}

		for _, replacement := range replacements {
			version, _, commit := parseGoModVersion(replace.New.Version)

			packages[replacement] = PackageDetails{
				Name:      replace.New.Path,
				Version:   strings.TrimPrefix(version, "v"),
				Ecosystem: GoEcosystem,
				CompareAs: GoEcosystem,
				Commit:    commit,
			}
		}
	}

	return pkgDetailsMapToSlice(deduplicatePackages(packages)), nil
}
