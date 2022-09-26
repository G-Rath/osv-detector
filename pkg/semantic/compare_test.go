package semantic_test

import (
	"github.com/g-rath/osv-detector/pkg/semantic"
	"math/big"
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
	a semantic.Version,
	b semantic.Version,
	expectedResult int,
) {
	t.Helper()

	if actualResult := a.Compare(b); actualResult != expectedResult {
		t.Errorf(
			"Expected %s to be %s %s, but it was %s",
			a.String(),
			compareWord(t, expectedResult),
			b.String(),
			compareWord(t, actualResult),
		)
	}
}

func buildlessVersion(build string, components ...int) semantic.Version {
	comps := make([]*big.Int, 0, len(components))

	for _, i := range components {
		comps = append(comps, big.NewInt(int64(i)))
	}

	return semantic.Version{Components: comps, Build: build}
}

func TestVersion_Compare_BasicEqual(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion(""),
		buildlessVersion(""),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1),
		buildlessVersion("", 1),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 2),
		buildlessVersion("", 1, 2),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 2, 3),
		buildlessVersion("", 1, 2, 3),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 2, 3, 4),
		buildlessVersion("", 1, 2, 3, 4),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 1, 0, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 1, 0, 0),
		0,
	)
}

func TestVersion_Compare_BasicWithBuildEqual(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc.1", 1),
		buildlessVersion("-rc.1", 1),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("beta2", 2, 0, 0),
		buildlessVersion("beta2", 2, 0, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(".v3", 1, 2, 3, 4, 5),
		buildlessVersion(".v3", 1, 2, 3, 4, 5),
		0,
	)
}

func TestVersion_Compare_BasicGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 2),
		buildlessVersion("", 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 2, 0),
		buildlessVersion("", 1, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 2, 0),
		buildlessVersion("", 0, 1, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 2, 0, 0),
		buildlessVersion("", 1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 1, 0),
		buildlessVersion("", 1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 1),
		buildlessVersion("", 1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 1, 1),
		buildlessVersion("", 1, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 4, 1, 1),
		buildlessVersion("", 1, 2, 3),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 2, 0),
		buildlessVersion("", 0, 0, 1),
		1,
	)
}

func TestVersion_Compare_BasicWithBuildGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc.2", 1),
		buildlessVersion("-rc.1", 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("beta22", 2, 0, 0),
		buildlessVersion("beta", 2, 0, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(".v20190411", 1, 2, 3, 4, 5),
		buildlessVersion(".v20190309", 1, 2, 3, 4, 5),
		1,
	)
}

func TestVersion_Compare_BasicLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 1),
		buildlessVersion("", 2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0),
		buildlessVersion("", 2, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 1, 0),
		buildlessVersion("", 0, 2, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 2, 0, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 1, 1, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 1, 0, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 0, 0),
		buildlessVersion("", 1, 1, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 1, 2, 3),
		buildlessVersion("", 4, 1, 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 0, 1),
		buildlessVersion("", 0, 2, 0),
		-1,
	)
}

func TestVersion_Compare_BasicWithBuildLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc.1", 1),
		buildlessVersion("-rc.2", 1),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("beta", 2, 0, 0),
		buildlessVersion("beta22", 2, 0, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(".v20190309", 1, 2, 3, 4, 5),
		buildlessVersion(".v20190411", 1, 2, 3, 4, 5),
		-1,
	)
}

func TestVersion_Compare_UnevenEquals(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 1),
		buildlessVersion("", 1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 1),
		buildlessVersion("", 0, 1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 0, 1),
		buildlessVersion("", 0, 0, 1, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 0, 0, 1),
		buildlessVersion("", 0, 0, 0, 1, 0),
		0,
	)
}

func TestVersion_Compare_UnevenWithBuildEqual(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc.1", 1),
		buildlessVersion("-rc.1", 1, 0, 0),
		0,
	)

	expectCompareResult(t,
		buildlessVersion("beta2", 0, 2, 0),
		buildlessVersion("beta2", 0, 2),
		0,
	)

	expectCompareResult(t,
		buildlessVersion(".v3", 1, 2, 3),
		buildlessVersion(".v3", 1, 2, 3, 0, 0),
		0,
	)
}

func TestVersion_Compare_UnevenGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 2, 2),
		buildlessVersion("", 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 2),
		buildlessVersion("", 0, 0, 0, 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 2, 0, 5),
		buildlessVersion("", 0, 1),
		1,
	)
}

func TestVersion_Compare_UnevenWithBuildGreaterThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc.1", 1),
		buildlessVersion("-rc", 1, 0),
		1,
	)

	expectCompareResult(t,
		buildlessVersion(".beta.5", 0, 2, 0),
		buildlessVersion(".alpha.2", 0, 2),
		1,
	)
}

func TestVersion_Compare_UnevenLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 1),
		buildlessVersion("", 2, 2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 0, 0, 1),
		buildlessVersion("", 2),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion("", 0, 1),
		buildlessVersion("", 0, 2, 0, 5),
		-1,
	)
}

func TestVersion_Compare_UnevenWithBuildLessThan(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("-rc", 1),
		buildlessVersion("-rc.1", 1, 0),
		-1,
	)

	expectCompareResult(t,
		buildlessVersion(".alpha.2", 0, 2),
		buildlessVersion(".beta.5", 0, 2, 0),
		-1,
	)
}

func TestVersion_Compare_MixedWithAndWithoutBuild(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		buildlessVersion("", 1),
		buildlessVersion("alpha", 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("rc.0", 2),
		buildlessVersion("", 1),
		1,
	)

	expectCompareResult(t,
		buildlessVersion("beta2", 1, 0),
		buildlessVersion("", 1),
		-1,
	)
}

// leading "v" is just cosmetic, and shouldn't change the comparing
func TestVersion_Compare_BasicWithLeadingV(t *testing.T) {
	t.Parallel()

	expectCompareResult(t,
		semantic.Version{LeadingV: false, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		semantic.Version{LeadingV: false, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		0,
	)

	expectCompareResult(t,
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		semantic.Version{LeadingV: false, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		0,
	)

	expectCompareResult(t,
		semantic.Version{LeadingV: false, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		0,
	)

	expectCompareResult(t,
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		0,
	)

	expectCompareResult(t,
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(2)}, Build: ""},
		semantic.Version{LeadingV: true, Components: []*big.Int{big.NewInt(1)}, Build: ""},
		1,
	)
}

func TestVersion_Compare_BasicWithBigComponents(t *testing.T) {
	t.Parallel()

	big1, _ := new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999", 10)
	big2, _ := new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999998", 10)

	expectCompareResult(t,
		semantic.Version{Components: []*big.Int{big1}},
		semantic.Version{Components: []*big.Int{big1}},
		0,
	)

	expectCompareResult(t,
		semantic.Version{Components: []*big.Int{big1}},
		semantic.Version{Components: []*big.Int{big2}},
		1,
	)

	expectCompareResult(t,
		semantic.Version{Components: []*big.Int{big1, big1, big1, big2}},
		semantic.Version{Components: []*big.Int{big1, big1, big1, big1}},
		-1,
	)
}
