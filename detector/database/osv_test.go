package database_test

import (
	"osv-detector/detector/database"
	"testing"
	"time"
)

func TestOSV_AffectsEcosystem(t *testing.T) {
	t.Parallel()

	type AffectsTest struct {
		Affected  []database.Affected
		Ecosystem database.Ecosystem
		Expected  bool
	}

	tests := []AffectsTest{
		{Affected: nil, Ecosystem: "Go", Expected: false},
		{Affected: nil, Ecosystem: "npm", Expected: false},
		{Affected: nil, Ecosystem: "PyPI", Expected: false},
		{Affected: nil, Ecosystem: "", Expected: false},
		{
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "crates.io"}},
				{Package: database.Package{Ecosystem: "npm"}},
				{Package: database.Package{Ecosystem: "PyPI"}},
			},
			Ecosystem: "Packagist",
			Expected:  false,
		},
		{
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "NuGet"}},
			},
			Ecosystem: "NuGet",
			Expected:  true,
		},
		{
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "npm"}},
				{Package: database.Package{Ecosystem: "npm"}},
			},
			Ecosystem: "npm",
			Expected:  true,
		},
	}

	for i, test := range tests {
		osv := database.OSV{
			ID:        "1",
			Published: time.Time{},
			Modified:  time.Time{},
			Details:   "This is an open source vulnerability!",
			Affected:  test.Affected,
		}

		if osv.AffectsEcosystem(test.Ecosystem) != test.Expected {
			t.Errorf(
				"Test #%d: Expected OSV to return %t but it returned %t",
				i,
				test.Expected,
				!test.Expected,
			)
		}
	}

	// test when the OSV doesn't have an "Affected"
	osv := database.OSV{
		ID:        "1",
		Published: time.Time{},
		Modified:  time.Time{},
		Details:   "This is an open source vulnerability!",
		Affected:  nil,
	}

	if osv.AffectsEcosystem("npm") {
		t.Errorf(
			"Expected OSV to report 'false' when it doesn't have an Affected, but it reported true!",
		)
	}
}
