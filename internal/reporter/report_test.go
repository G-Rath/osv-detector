package reporter_test

import (
	"osv-detector/internal"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
	"osv-detector/internal/reporter"
	"strings"
	"testing"
)

func TestReport_HasKnownVulnerabilities(t *testing.T) {
	t.Parallel()

	packageDetails := internal.PackageDetails{
		Name:      "addr2line",
		Version:   "0.15.2",
		Ecosystem: lockfile.CargoEcosystem,
	}
	type fields struct {
		Lockfile lockfile.Lockfile
		Packages []reporter.PackageDetailsWithVulnerabilities
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "no packages",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{},
			},
			want: false,
		},
		{
			name: "no vulnerabilities",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{},
					},
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{},
					},
				},
			},
			want: false,
		},
		{
			name: "one package with one vulnerability",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{{ID: "1"}},
					},
				},
			},
			want: true,
		},
		{
			name: "multiple packages, with one with a vulnerability",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{},
					},
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{{ID: "1"}},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := reporter.Report{
				Lockfile: tt.fields.Lockfile,
				Packages: tt.fields.Packages,
			}
			if got := r.HasKnownVulnerabilities(); got != tt.want {
				t.Errorf("HasKnownVulnerabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport_ToString_NoVulnerabilities(t *testing.T) {
	t.Parallel()

	msg := "no known vulnerabilities found"

	r := reporter.Report{}

	if actual := r.ToString(); !strings.Contains(actual, msg) {
		t.Errorf("Expected \"%s\" to contain \"%s\" but it did not", actual, msg)
	}
}

func TestReport_ToString_OneVulnerability(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  my-package@1.2.3 is affected by the following vulnerabilities:",
		"    GHSA-1: This is a vulnerability! (https://github.com/advisories/GHSA-1)",
		"",
		"  1 known vulnerability found in /path/to/my/lock",
		"",
	}, "\n")

	r := reporter.Report{
		Lockfile: lockfile.Lockfile{FilePath: "/path/to/my/lock"},
		Packages: []reporter.PackageDetailsWithVulnerabilities{
			{
				PackageDetails: internal.PackageDetails{
					Name:      "my-package",
					Version:   "1.2.3",
					Ecosystem: lockfile.BundlerEcosystem,
				},
				Vulnerabilities: []database.OSV{
					{
						ID:      "GHSA-1",
						Summary: "This is a vulnerability!",
					},
				},
			},
		},
	}

	if actual := r.ToString(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestReport_ToString_MultipleVulnerabilities(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  my-package@1.2.3 is affected by the following vulnerabilities:",
		"    GHSA-1: This is a vulnerability! (https://github.com/advisories/GHSA-1)",
		"  their-package@4.5.6 is affected by the following vulnerabilities:",
		"    GHSA-2: This is another vulnerability! (https://github.com/advisories/GHSA-2)",
		"",
		"  2 known vulnerabilities found in /path/to/my/lock",
		"",
	}, "\n")

	r := reporter.Report{
		Lockfile: lockfile.Lockfile{FilePath: "/path/to/my/lock"},
		Packages: []reporter.PackageDetailsWithVulnerabilities{
			{
				PackageDetails: internal.PackageDetails{
					Name:      "my-package",
					Version:   "1.2.3",
					Ecosystem: lockfile.BundlerEcosystem,
				},
				Vulnerabilities: []database.OSV{
					{
						ID:      "GHSA-1",
						Summary: "This is a vulnerability!",
					},
				},
			},
			{
				PackageDetails: internal.PackageDetails{
					Name:      "middle-package",
					Version:   "1.2.0",
					Ecosystem: lockfile.BundlerEcosystem,
				},
				Vulnerabilities: []database.OSV{},
			},
			{
				PackageDetails: internal.PackageDetails{
					Name:      "their-package",
					Version:   "4.5.6",
					Ecosystem: lockfile.BundlerEcosystem,
				},
				Vulnerabilities: []database.OSV{
					{
						ID:      "GHSA-2",
						Summary: "This is another vulnerability!",
					},
				},
			},
		},
	}

	if actual := r.ToString(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}
