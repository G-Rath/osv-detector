package semantic_test

import (
	"bufio"
	"fmt"
	"os"
	"osv-detector/internal/semantic"
	"strings"
	"testing"
)

func versionsEqual(expectedVersion semantic.Version, actualVersion semantic.Version) bool {
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

func explainVersion(version semantic.Version) string {
	str := "{ "

	for i, component := range version.Components {
		str += fmt.Sprintf("%d: %d, ", i, component)
	}

	str += fmt.Sprintf("Build: %s", version.Build)
	str += " }"

	return str
}

func notExpectedMessage(version string, expectedVersion semantic.Version, actualVersion semantic.Version) string {
	str := fmt.Sprintf("'%s' was not parsed as expected:", version)

	str += fmt.Sprintf("\n  Expected: %s", explainVersion(expectedVersion))
	str += fmt.Sprintf("\n  Actual  : %s", explainVersion(actualVersion))

	return str
}

func expectParsedVersionToMatchOriginalString(t *testing.T, str string) semantic.Version {
	t.Helper()

	actualVersion := semantic.Parse(str)

	if actualVersion.ToString() != str {
		t.Errorf(
			"Parsed version as a string did not equal original: %s != %s",
			actualVersion.ToString(),
			str,
		)
	}

	return actualVersion
}

func expectParsedAsVersion(t *testing.T, str string, expectedVersion semantic.Version) {
	t.Helper()

	actualVersion := expectParsedVersionToMatchOriginalString(t, str)

	if !versionsEqual(expectedVersion, actualVersion) {
		t.Errorf(notExpectedMessage(str, expectedVersion, actualVersion))
	}
}

func expectParsedVersionToMatchString(
	t *testing.T,
	str string,
	expectedString string,
	expectedVersion semantic.Version,
) {
	t.Helper()

	actualVersion := semantic.Parse(str)

	if actualVersion.ToString() != expectedString {
		t.Errorf(
			"Parsed version as a string did not equal expected: %s != %s",
			actualVersion.ToString(),
			expectedString,
		)
	}

	if !versionsEqual(expectedVersion, actualVersion) {
		t.Errorf(notExpectedMessage(str, expectedVersion, actualVersion))
	}
}

func TestParse_Standard(t *testing.T) {
	t.Parallel()

	expectParsedAsVersion(t, "0.0.0.0", semantic.Version{
		Components: []int{0, 0, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.0.0.0", semantic.Version{
		Components: []int{1, 0, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.0.0", semantic.Version{
		Components: []int{1, 2, 0, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.0", semantic.Version{
		Components: []int{1, 2, 3, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.4", semantic.Version{
		Components: []int{1, 2, 3, 4},
		Build:      "",
	})

	expectParsedAsVersion(t, "9.2.55826.0", semantic.Version{
		Components: []int{9, 2, 55826, 0},
		Build:      "",
	})

	expectParsedAsVersion(t, "3.2.22.3", semantic.Version{
		Components: []int{3, 2, 22, 3},
		Build:      "",
	})
}

func TestParse_Omitted(t *testing.T) {
	t.Parallel()

	expectParsedAsVersion(t, "1", semantic.Version{
		Components: []int{1},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2", semantic.Version{
		Components: []int{1, 2},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3", semantic.Version{
		Components: []int{1, 2, 3},
		Build:      "",
	})

	expectParsedAsVersion(t, "1.2.3.", semantic.Version{
		Components: []int{1, 2, 3},
		Build:      ".",
	})
}

func TestParse_WithBuildString(t *testing.T) {
	t.Parallel()

	expectParsedAsVersion(t, "10.0.0.beta1", semantic.Version{
		Components: []int{10, 0, 0},
		Build:      ".beta1",
	})

	expectParsedAsVersion(t, "1.0.0a20", semantic.Version{
		Components: []int{1, 0, 0},
		Build:      "a20",
	})

	expectParsedAsVersion(t, "9.0.0.pre1", semantic.Version{
		Components: []int{9, 0, 0},
		Build:      ".pre1",
	})

	expectParsedAsVersion(t, "9.4.16.v20190411", semantic.Version{
		Components: []int{9, 4, 16},
		Build:      ".v20190411",
	})

	expectParsedAsVersion(t, "0.3.0-beta.83", semantic.Version{
		Components: []int{0, 3, 0},
		Build:      "-beta.83",
	})

	expectParsedAsVersion(t, "3.0.0-beta.17.5", semantic.Version{
		Components: []int{3, 0, 0},
		Build:      "-beta.17.5",
	})

	expectParsedAsVersion(t, "4.0.0-milestone3", semantic.Version{
		Components: []int{4, 0, 0},
		Build:      "-milestone3",
	})

	expectParsedAsVersion(t, "13.6RC1", semantic.Version{
		Components: []int{13, 6},
		Build:      "RC1",
	})
}

func TestParse_MassParsing(t *testing.T) {
	t.Parallel()

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

func TestParse_NoComponents(t *testing.T) {
	t.Parallel()

	expectParsedVersionToMatchOriginalString(t, "hello world!")
}

func TestParse_LeadingZerosAndDateLike(t *testing.T) {
	t.Parallel()

	expectParsedVersionToMatchString(t, "20.04.0", "20.4.0", semantic.Version{
		Components: []int{20, 4, 0},
		Build:      "",
	})

	expectParsedVersionToMatchString(t, "4.3.04", "4.3.4", semantic.Version{
		Components: []int{4, 3, 4},
		Build:      "",
	})
}

// some versions look like they might be dates, which currently is not supported
// technically because we're using ints so their leading zeros are discarded,
// but practically we can't really know for sure if a version string should be
// treated as a date, so for now we're just treating them as versions
//
// todo: look into this more, and confirm if these versions are actually
//  meant to be dates, and are expected to be compared as such
func TestParse_DateLike(t *testing.T) {
	t.Parallel()

	expectParsedVersionToMatchString(t, "20.04.0", "20.4.0", semantic.Version{
		Components: []int{20, 4, 0},
		Build:      "",
	})

	expectParsedVersionToMatchString(t, "4.3.04alpha01", "4.3.4alpha01", semantic.Version{
		Components: []int{4, 3, 4},
		Build:      "alpha01",
	})

	expectParsedVersionToMatchString(t, "2019.03.6.1", "2019.3.6.1", semantic.Version{
		Components: []int{2019, 3, 6, 1},
		Build:      "",
	})

	expectParsedVersionToMatchString(t, "19.04.15", "19.4.15", semantic.Version{
		Components: []int{19, 4, 15},
		Build:      "",
	})

	expectParsedVersionToMatchString(t, "20.04.13", "20.4.13", semantic.Version{
		Components: []int{20, 4, 13},
		Build:      "",
	})

	expectParsedVersionToMatchString(t, "2019.11.09", "2019.11.9", semantic.Version{
		Components: []int{2019, 11, 9},
		Build:      "",
	})
}
