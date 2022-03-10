package database_test

import (
	"osv-detector/detector"
	"osv-detector/detector/database"
	"osv-detector/detector/parsers"
	"testing"
	"time"
)

func expectIsAffected(t *testing.T, osv database.OSV, version string, expectAffected bool) {
	t.Helper()

	pkg := detector.PackageDetails{Name: "my-package", Version: version, Ecosystem: parsers.NpmEcosystem}

	if osv.IsAffected(pkg) != expectAffected {
		if expectAffected {
			t.Errorf("Expected OSV to affect package version %s but it did not", version)
		} else {
			t.Errorf("Expected OSV not to affect package version %s but it did", version)
		}
	}
}

func buildOSVWithAffected(affected ...database.Affected) database.OSV {
	return database.OSV{
		ID:        "1",
		Published: time.Time{},
		Modified:  time.Time{},
		Details:   "This is an open source vulnerability!",
		Affected:  affected,
	}
}

func buildEcosystemAffectsRange(events ...database.RangeEvent) database.AffectsRange {
	return database.AffectsRange{Type: database.TypeEcosystem, Events: events}
}

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

func TestOSV_IsAffected_AffectsWithEcosystem_DifferentEcosystem(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: parsers.PipEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(database.RangeEvent{Introduced: "0"}),
			},
		},
	)

	for _, v := range []string{"1.0.0", "1.1.1", "2.0.0"} {
		expectIsAffected(t, osv, v, false)
	}
}

func TestOSV_IsAffected_AffectsWithEcosystem_SingleAffected(t *testing.T) {
	t.Parallel()

	var osv database.OSV

	// "Introduced: 0" means everything is affected
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: parsers.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(database.RangeEvent{Introduced: "0"}),
			},
		},
	)

	for _, v := range []string{"1.0.0", "1.1.1", "2.0.0"} {
		expectIsAffected(t, osv, v, true)
	}

	// "Fixed: 1" means all versions after this are not vulnerable
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: parsers.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.1.0", "0.0.0.1", "1.0.0-rc"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"1.0.0", "1.1.0", "2.0.0"} {
		expectIsAffected(t, osv, v, false)
	}

	// multiple fixes and introduced
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: parsers.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "3.2.0"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.1.0", "0.0.0.1", "1.0.0-rc"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"1.0.0", "1.1.0", "2.0.0rc2", "2.0.1"} {
		expectIsAffected(t, osv, v, false)
	}

	for _, v := range []string{"2.1.1", "2.3.4", "3.0.0", "3.0.0-rc"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"3.2.0", "3.2.1", "4.0.0"} {
		expectIsAffected(t, osv, v, false)
	}
}

func TestOSV_IsAffected_AffectsWithEcosystem_MultipleAffected(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: parsers.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: parsers.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "3.2.0"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.1.0", "0.0.0.1", "1.0.0-rc"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"1.0.0", "1.1.0", "2.0.0rc2", "2.0.1"} {
		expectIsAffected(t, osv, v, false)
	}

	for _, v := range []string{"2.1.1", "2.3.4", "3.0.0", "3.0.0-rc"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"3.2.0", "3.2.1", "4.0.0"} {
		expectIsAffected(t, osv, v, false)
	}
}
