package semver_test

import (
	"bufio"
	"fmt"
	"os"
	"osv-detector/detector/semver"
	"strings"
	"testing"
)

func versionsEqual(expectedVersion semver.Version, actualVersion semver.Version) bool {
	if expectedVersion.Build != actualVersion.Build {
		return false
	}

	if len(expectedVersion.Components) != len(actualVersion.Components) {
		return false
	}

	for i := range expectedVersion.Components {
		if expectedVersion.Components[i] != actualVersion.Components[i] {
			return false
		}
	}

	return true
}

func explainVersion(version semver.Version) string {
	str := "{ "

	for i, component := range version.Components {
		str += fmt.Sprintf("%d: %d, ", i, component)
	}

	str += fmt.Sprintf("Build: %s", version.Build)
	str += " }"

	return str
}

func notExpectedMessage(version string, expectedVersion semver.Version, actualVersion semver.Version) string {
	str := fmt.Sprintf("'%s' was not parsed as expected:", version)

	str += fmt.Sprintf("\n  Expected: %s", explainVersion(expectedVersion))
	str += fmt.Sprintf("\n  Actual  : %s", explainVersion(actualVersion))

	return str
}

func expectParsedVersionToMatchOriginalString(t *testing.T, str string) semver.Version {
	t.Helper()

	actualVersion := semver.Parse(str)

	if actualVersion.ToString() != str {
		t.Errorf(
			"Parsed version as a string did not equal original: %s != %s",
			actualVersion.ToString(),
			str,
		)
	}

	return actualVersion
}

func expectParsedAsVersion(t *testing.T, str string, expectedVersion semver.Version) {
	t.Helper()

	actualVersion := expectParsedVersionToMatchOriginalString(t, str)

	if !versionsEqual(expectedVersion, actualVersion) {
		t.Errorf(notExpectedMessage(str, expectedVersion, actualVersion))
	}
}

func TestParse_Standard(t *testing.T) {
	expectParsedAsVersion(t, "0.0.0.0", semver.Version{
		Components: []int{0, 0, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.0.0.0", semver.Version{
		Components: []int{1, 0, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.0.0", semver.Version{
		Components: []int{1, 2, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.0", semver.Version{
		Components: []int{1, 2, 3, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.4", semver.Version{
		Components: []int{1, 2, 3, 4},
		Build:      "",
	})

	expectParsedAsVersion(t, "9.2.55826.0", semver.Version{
		Components: []int{9, 2, 55826, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "3.2.22.3", semver.Version{
		Components: []int{3, 2, 22, 3},
		Build:      "",
	})
}

func TestParse_Omitted(t *testing.T) {
	expectParsedAsVersion(t, "1", semver.Version{
		Components: []int{1},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2", semver.Version{
		Components: []int{1, 2},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3", semver.Version{
		Components: []int{1, 2, 3},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.", semver.Version{
		Components: []int{1, 2, 3},
		Build:      ".",
	})
}

func TestParse_WithBuildString(t *testing.T) {
	expectParsedAsVersion(t, "10.0.0.beta1", semver.Version{
		Components: []int{10, 0, 0},
		Build:      ".beta1",
	})

	expectParsedAsVersion(t, "1.0.0a20", semver.Version{
		Components: []int{1, 0, 0},
		Build:      "a20",
	})

	expectParsedAsVersion(t, "9.0.0.pre1", semver.Version{
		Components: []int{9, 0, 0},
		Build:      ".pre1",
	})

	expectParsedAsVersion(t, "9.4.16.v20190411", semver.Version{
		Components: []int{9, 4, 16},
		Build:      ".v20190411",
	})

	expectParsedAsVersion(t, "0.3.0-beta.83", semver.Version{
		Components: []int{0, 3, 0},
		Build:      "-beta.83",
	})

	expectParsedAsVersion(t, "3.0.0-beta.17.5", semver.Version{
		Components: []int{3, 0, 0},
		Build:      "-beta.17.5",
	})

	expectParsedAsVersion(t, "4.0.0-milestone3", semver.Version{
		Components: []int{4, 0, 0},
		Build:      "-milestone3",
	})

	expectParsedAsVersion(t, "13.6RC1", semver.Version{
		Components: []int{13, 6},
		Build:      "RC1",
	})
}

func TestParse_MassParsing(t *testing.T) {
	file, err := os.Open("fixtures/all-versions.txt")
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// treat "#" as comments because they're not supported (yet?)
		if strings.HasPrefix(line, "# ") {
			continue
		}

		expectParsedVersionToMatchOriginalString(t, line)
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
