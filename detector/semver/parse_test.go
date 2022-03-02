package semver_test

import (
	"fmt"
	"osv-detector/detector/semver"
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

func expectParsedAsVersion(t *testing.T, str string, expectedVersion semver.Version) {
	t.Helper()

	actualVersion := semver.Parse(str)

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
}
