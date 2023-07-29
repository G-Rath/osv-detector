package database_test

import (
	"testing"
	"time"

	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/models"
)

func expectOSVDescription(t *testing.T, expected string, osv database.OSV) {
	t.Helper()

	if actual := osv.Describe(); actual != expected {
		t.Errorf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func expectIsAffected(t *testing.T, osv database.OSV, version string, expectAffected bool) {
	t.Helper()

	pkg := internal.PackageDetails{
		Name:      "my-package",
		Version:   version,
		Ecosystem: models.EcosystemNPM,
		CompareAs: models.EcosystemNPM,
	}

	if osv.IsAffected(pkg) != expectAffected {
		if version == "" {
			version = "<empty>"
		}

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

func buildSemverAffectsRange(events ...database.RangeEvent) database.AffectsRange {
	return database.AffectsRange{Type: database.TypeSemver, Events: events}
}

func TestPackage_NormalizedName(t *testing.T) {
	t.Parallel()

	database.Package{
		Name:      "",
		Ecosystem: "",
	}.NormalizedName()
}

func TestPackage_NormalizedName_PipEcosystem(t *testing.T) {
	t.Parallel()

	x := [][]string{
		{"Pillow", "pillow"},
		{"privacyIDEA", "privacyidea"},
		{"Products.GenericSetup", "products-genericsetup"},
	}

	for _, strings := range x {
		name := database.Package{Name: strings[0], Ecosystem: models.EcosystemPyPI}.NormalizedName()

		if name != strings[1] {
			t.Errorf(
				"Expected package named %s to be normalized to %s but was normalized to %s instead",
				strings[0],
				strings[1],
				name,
			)
		}
	}
}

func TestPackage_NormalizedName_NotPipEcosystem(t *testing.T) {
	t.Parallel()

	x := []string{
		"Proto",
		"Pillow",
		"privacyIDEA",
		"Products.GenericSetup",
	}

	for _, na := range x {
		name := database.Package{Name: na, Ecosystem: models.EcosystemNPM}.NormalizedName()

		if name != na {
			t.Errorf(
				"Expected package named %s to be unchanged, but it was normalized to %s",
				na,
				name,
			)
		}
	}
}

func TestOSV_AffectsEcosystem(t *testing.T) {
	t.Parallel()

	type AffectsTest struct {
		Affected  []database.Affected
		Ecosystem models.Ecosystem
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
			Package: database.Package{Ecosystem: models.EcosystemPyPI, Name: "my-package"},
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
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(database.RangeEvent{Introduced: "0"}),
			},
		},
	)

	for _, v := range []string{"1.0.0", "1.1.1", "2.0.0"} {
		expectIsAffected(t, osv, v, true)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// "Fixed: 1" means all versions after this are not vulnerable
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// multiple fixes and introduced
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// multiple fixes and introduced, shuffled
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "3.2.0"},
					database.RangeEvent{Fixed: "1"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// "LastAffected: 1" means all versions after this are not vulnerable
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{LastAffected: "1"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.1.0", "0.0.0.1", "1.0.0-rc", "1.0.0"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"1.0.1", "1.1.0", "2.0.0"} {
		expectIsAffected(t, osv, v, false)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// mix of fixes, last_known_affected, and introduced
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{LastAffected: "3.1.9"},
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

	for _, v := range []string{"2.1.1", "2.3.4", "3.0.0", "3.0.0-rc", "3.1.9"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"3.2.0", "3.2.1", "4.0.0"} {
		expectIsAffected(t, osv, v, false)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// mix of fixes, last_known_affected, and introduced, shuffled
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{LastAffected: "3.1.9"},
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

	for _, v := range []string{"2.1.1", "2.3.4", "3.0.0", "3.0.0-rc", "3.1.9"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"3.2.0", "3.2.1", "4.0.0"} {
		expectIsAffected(t, osv, v, false)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)
}

func TestOSV_IsAffected_AffectsWithEcosystem_MultipleAffected(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "3.2.0"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "3.3.0"},
					database.RangeEvent{LastAffected: "3.5.0"},
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

	for _, v := range []string{"3.3.1", "3.4.5"} {
		expectIsAffected(t, osv, v, true)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// shuffled
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{Introduced: "0"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Fixed: "3.2.0"},
					database.RangeEvent{Introduced: "2.1.0"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{LastAffected: "3.5.0"},
					database.RangeEvent{Introduced: "3.3.0"},
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

	for _, v := range []string{"3.3.1", "3.4.5"} {
		expectIsAffected(t, osv, v, true)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// zeros with build strings
	osv = buildOSVWithAffected(
		database.Affected{
			// golang.org/x/sys
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Fixed: "0.0.0-20220412211240-33da011f77ad"},
					database.RangeEvent{Introduced: "0"},
				),
			},
		},
		database.Affected{
			// golang.org/x/net
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0.0.0-20180925071336-cf3bd585ca2a"},
					database.RangeEvent{Fixed: "0"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.14.0"} {
		expectIsAffected(t, osv, v, false)
	}

	for _, v := range []string{"0.0.0-20180925071336-cf3bd585ca2a"} {
		expectIsAffected(t, osv, v, true)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)
}

func TestOSV_IsAffected_AffectsWithEcosystem_PipNamesAreNormalised(t *testing.T) {
	t.Parallel()

	var osv database.OSV
	var pkg internal.PackageDetails

	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemPyPI, Name: "Pillow"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
	)

	pkg = internal.PackageDetails{
		Name:      "pillow",
		Version:   "0.5",
		Ecosystem: models.EcosystemPyPI,
		CompareAs: models.EcosystemPyPI,
	}

	if !osv.IsAffected(pkg) {
		t.Errorf("Expected OSV to normalize names of pip packages, but did not")
	}

	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "Pillow"},
			Ranges: []database.AffectsRange{
				buildEcosystemAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
	)

	pkg = internal.PackageDetails{
		Name:      "pillow",
		Version:   "0.5",
		Ecosystem: models.EcosystemNPM,
		CompareAs: models.EcosystemNPM,
	}

	if osv.IsAffected(pkg) {
		t.Errorf("Expected OSV not to normalize names of non-pip packages, but it did")
	}
}

func TestOSV_IsAffected_AffectsWithSemver_DifferentEcosystem(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemPyPI, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(database.RangeEvent{Introduced: "0"}),
			},
		},
	)

	for _, v := range []string{"1.0.0", "1.1.1", "2.0.0"} {
		expectIsAffected(t, osv, v, false)
	}
}

func TestOSV_IsAffected_AffectsWithSemver_SingleAffected(t *testing.T) {
	t.Parallel()

	var osv database.OSV

	// "Introduced: 0" means everything is affected
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(database.RangeEvent{Introduced: "0"}),
			},
		},
	)

	for _, v := range []string{"v1.0.0", "v1.1.1", "v2.0.0"} {
		expectIsAffected(t, osv, v, true)
	}

	// "Fixed: 1" means all versions after this are not vulnerable
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1.0.0"},
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
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// multiple fixes and introduced, shuffled
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{Fixed: "3.2.0"},
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Introduced: "2.1.0"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// "LastAffected: 1" means all versions after this are not vulnerable
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{LastAffected: "1.0.0"},
				),
			},
		},
	)

	for _, v := range []string{"0.0.0", "0.1.0", "0.0.0.1", "1.0.0-rc", "1.0.0"} {
		expectIsAffected(t, osv, v, true)
	}

	for _, v := range []string{"1.0.1", "1.1.0", "2.0.0"} {
		expectIsAffected(t, osv, v, false)
	}

	// mix of fixes, last_known_affected, and introduced
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{LastAffected: "3.1.9"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)

	// mix of fixes, last_known_affected, and introduced, shuffled
	osv = buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: lockfile.NpmEcosystem, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{LastAffected: "3.1.9"},
					database.RangeEvent{Fixed: "1"},
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

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)
}

func TestOSV_IsAffected_AffectsWithSemver_MultipleAffected(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "0"},
					database.RangeEvent{Fixed: "1"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "2.1.0"},
					database.RangeEvent{Fixed: "3.2.0"},
				),
			},
		},
		database.Affected{
			Package: database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Ranges: []database.AffectsRange{
				buildSemverAffectsRange(
					database.RangeEvent{Introduced: "3.3.0"},
					database.RangeEvent{LastAffected: "3.5.0"},
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

	for _, v := range []string{"3.3.1", "3.4.5", "3.5.0"} {
		expectIsAffected(t, osv, v, true)
	}

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)
}

func TestOSV_IsAffected_OnlyVersions(t *testing.T) {
	t.Parallel()

	osv := buildOSVWithAffected(
		database.Affected{
			Package:  database.Package{Ecosystem: models.EcosystemNPM, Name: "my-package"},
			Versions: []string{"1.0.0"},
		},
	)

	expectIsAffected(t, osv, "0.0.0", false)
	expectIsAffected(t, osv, "1.0.0", true)
	expectIsAffected(t, osv, "1.0.0-beta1", false)
	expectIsAffected(t, osv, "1.1.0", false)

	// an empty version should always be treated as affected
	expectIsAffected(t, osv, "", true)
}

func TestOSV_Describe_Text(t *testing.T) {
	t.Parallel()

	// use the summary of the advisory if present
	expectOSVDescription(t, "This is a vulnerability!", database.OSV{Summary: "This is a vulnerability!"})

	// otherwise, use the details of the advisory
	expectOSVDescription(t, "It's very bad!", database.OSV{Details: "It's very bad!"})

	// prefer the summary over details, as the former should be more concise
	expectOSVDescription(t,
		"This is a vulnerability!",
		database.OSV{Summary: "This is a vulnerability!", Details: "It's very bad!"},
	)

	// advisories from GitHub should have a link included with their description
	expectOSVDescription(t,
		"This is a vulnerability! (https://github.com/advisories/GHSA-1)",
		database.OSV{Summary: "This is a vulnerability!", ID: "GHSA-1"},
	)

	// if none of those are present, say that there are no details available
	expectOSVDescription(t, "(no details available)", database.OSV{})
}

func TestOSV_Describe_Truncation(t *testing.T) {
	t.Parallel()

	// long details text should be truncated after 80 characters, based on words
	expectOSVDescription(t,
		"sqlparse is a non-validating SQL parser module for Python. In sqlparse versions...",
		database.OSV{
			Details: "sqlparse is a non-validating SQL parser module for Python. In sqlparse versions 0.4.0 and 0.4.1 there is a regular Expression Denial of Service in sqlparse vulnerability.",
		},
	)

	// if a link is present, that isn't included in the truncating
	expectOSVDescription(t,
		"sqlparse is a non-validating SQL parser module for Python. In sqlparse versions... (https://github.com/advisories/GHSA-1)",
		database.OSV{
			Details: "sqlparse is a non-validating SQL parser module for Python. In sqlparse versions 0.4.0 and 0.4.1 there is a regular Expression Denial of Service in sqlparse vulnerability.",
			ID:      "GHSA-1",
		},
	)

	// long continuous text without any spaces before the limit should be forcefully truncated
	expectOSVDescription(t,
		"nannannannannannannannannannannannannannannannannannannannannannannannannannanna...",
		database.OSV{
			Details: "nannannannannannannannannannannannannannannannannannannannannannannannannannannannanan batman!",
		},
	)

	// truncation shouldn't be applied to the summary text (as it should be short already)
	expectOSVDescription(t,
		"This is a vulnerability! It's very serious, and should not be taken lightly or bad things could happen!",
		database.OSV{
			Summary: "This is a vulnerability! It's very serious, and should not be taken lightly or bad things could happen!",
		},
	)
}

func TestVersions_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vs   database.Versions
		want string
	}{
		{
			name: "",
			vs:   nil,
			want: "[]",
		},
		{
			name: "",
			vs:   database.Versions(nil),
			want: "[]",
		},
		{
			name: "",
			vs:   database.Versions{"1.0.0"},
			want: "[\"1.0.0\"]",
		},
		{
			name: "",
			vs:   database.Versions{"1.0.0", "1.2.3", "4.5.6"},
			want: "[\"1.0.0\",\"1.2.3\",\"4.5.6\"]",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.vs.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)

				return
			}

			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", gotStr, tt.want)
			}
		})
	}
}
