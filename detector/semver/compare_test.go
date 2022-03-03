package semver_test

import (
	"osv-detector/detector/semver"
	"testing"
)

func compareWord(t *testing.T, result int) string {
	t.Helper()

	switch result {
	case 1:
		return "greater than"
	case 0:
		return "equal to"
	case -1:
		return "less than"
	default:
		t.Fatalf("Unexpected compare result: %d\n", result)

		return ""
	}
}

func expectCompareResult(
	t *testing.T,
	a semver.Version,
	b semver.Version,
	expectedResult int,
) {
	t.Helper()

	if actualResult := a.Compare(b); actualResult != expectedResult {
		t.Errorf(
			"Expected %s to be %s %s, but it was %s",
			a.ToString(),
			compareWord(t, expectedResult),
			b.ToString(),
			compareWord(t, actualResult),
		)
	}
}

func buildlessVersion(components ...int) semver.Version {
	return semver.Version{Components: components}
}

func TestVersion_Compare_BasicEqual(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(),
		buildlessVersion(),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1),
		buildlessVersion(1),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1, 2),
		buildlessVersion(1, 2),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1, 2, 3),
		buildlessVersion(1, 2, 3),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1, 2, 3, 4),
		buildlessVersion(1, 2, 3, 4),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(1, 0, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(1, 0, 0),
		0,
	)
}

func TestVersion_Compare_BasicGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(2),
		buildlessVersion(1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(2, 0),
		buildlessVersion(1, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 2, 0),
		buildlessVersion(0, 1, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(2, 0, 0),
		buildlessVersion(1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 1, 0),
		buildlessVersion(1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 1),
		buildlessVersion(1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 1, 1),
		buildlessVersion(1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(4, 1, 1),
		buildlessVersion(1, 2, 3),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 2, 0),
		buildlessVersion(0, 0, 1),
		1,
	)
}

func TestVersion_Compare_BasicLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(1),
		buildlessVersion(2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0),
		buildlessVersion(2, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 1, 0),
		buildlessVersion(0, 2, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(2, 0, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(1, 1, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(1, 0, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 0, 0),
		buildlessVersion(1, 1, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(1, 2, 3),
		buildlessVersion(4, 1, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 0, 1),
		buildlessVersion(0, 2, 0),
		-1,
	)
}

func TestVersion_Compare_UnevenEquals(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(1),
		buildlessVersion(1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(0, 1),
		buildlessVersion(0, 1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(0, 0, 1),
		buildlessVersion(0, 0, 1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(0, 0, 0, 1),
		buildlessVersion(0, 0, 0, 1, 0),
		0,
	)
}

func TestVersion_Compare_UnevenGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(2, 2),
		buildlessVersion(1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(2),
		buildlessVersion(0, 0, 0, 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 2, 0, 5),
		buildlessVersion(0, 1),
		1,
	)
}

func TestVersion_Compare_UnevenLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(1),
		buildlessVersion(2, 2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 0, 0, 1),
		buildlessVersion(2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(0, 1),
		buildlessVersion(0, 2, 0, 5),
		-1,
	)
}
