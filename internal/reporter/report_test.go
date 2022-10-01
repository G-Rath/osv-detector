package reporter_test

import (
	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/internal/reporter"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/lockfile"
	"strings"
	"testing"
)

func TestReport_HasKnownVulnerabilities(t *testing.T) {
	t.Parallel()

	packageDetails := internal.PackageDetails{
		Name:      "addr2line",
		Version:   "0.15.2",
		Ecosystem: lockfile.CargoEcosystem,
		CompareAs: lockfile.CargoEcosystem,
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

func TestReport_HasIgnoredVulnerabilities(t *testing.T) {
	t.Parallel()

	packageDetails := internal.PackageDetails{
		Name:      "addr2line",
		Version:   "0.15.2",
		Ecosystem: lockfile.CargoEcosystem,
		CompareAs: lockfile.CargoEcosystem,
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
			want: false,
		},
		{
			name: "one package with one ignored vulnerability",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{},
						Ignored:         []database.OSV{{ID: "1"}},
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
			want: false,
		},
		{
			name: "multiple packages, with one with a ignored vulnerability",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{},
					},
					{
						PackageDetails: packageDetails,
						Ignored:        []database.OSV{{ID: "1"}},
					},
				},
			},
			want: true,
		},
		{
			name: "multiple packages, with one with a vulnerability and one ignored",
			fields: fields{
				Lockfile: lockfile.Lockfile{},
				Packages: []reporter.PackageDetailsWithVulnerabilities{
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{{ID: "1"}},
					},
					{
						PackageDetails:  packageDetails,
						Vulnerabilities: []database.OSV{{ID: "1"}},
						Ignored:         []database.OSV{{ID: "2"}},
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
			if got := r.HasIgnoredVulnerabilities(); got != tt.want {
				t.Errorf("HasIgnoredVulnerabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport_String_NoVulnerabilities(t *testing.T) {
	t.Parallel()

	msg := "no known vulnerabilities found"

	r := reporter.Report{}

	if actual := r.String(); !strings.Contains(actual, msg) {
		t.Errorf("Expected \"%s\" to contain \"%s\" but it did not", actual, msg)
	}
}

func TestReport_String_OneVulnerability(t *testing.T) {
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
					CompareAs: lockfile.BundlerEcosystem,
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

	if actual := r.String(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestReport_String_MultipleVulnerabilities(t *testing.T) {
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
					CompareAs: lockfile.BundlerEcosystem,
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
					CompareAs: lockfile.BundlerEcosystem,
				},
				Vulnerabilities: []database.OSV{},
			},
			{
				PackageDetails: internal.PackageDetails{
					Name:      "their-package",
					Version:   "4.5.6",
					Ecosystem: lockfile.BundlerEcosystem,
					CompareAs: lockfile.BundlerEcosystem,
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

	if actual := r.String(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestReport_String_AllIgnoredVulnerabilities(t *testing.T) {
	t.Parallel()

	msg := "no new vulnerabilities found (2 were ignored)"

	r := reporter.Report{
		Lockfile: lockfile.Lockfile{FilePath: "/path/to/my/lock"},
		Packages: []reporter.PackageDetailsWithVulnerabilities{
			{
				PackageDetails: internal.PackageDetails{
					Name:      "my-package",
					Version:   "1.2.3",
					Ecosystem: lockfile.BundlerEcosystem,
					CompareAs: lockfile.BundlerEcosystem,
				},
				Ignored: []database.OSV{
					{
						ID:      "GHSA-1",
						Summary: "This is a vulnerability!",
					},
				},
			},
			{
				PackageDetails: internal.PackageDetails{
					Name:      "their-package",
					Version:   "4.5.6",
					Ecosystem: lockfile.BundlerEcosystem,
					CompareAs: lockfile.BundlerEcosystem,
				},
				Ignored: []database.OSV{
					{
						ID:      "GHSA-2",
						Summary: "This is another vulnerability!",
					},
				},
			},
		},
	}

	if actual := r.String(); !strings.Contains(actual, msg) {
		t.Errorf("Expected \"%s\" to contain \"%s\" but it did not", actual, msg)
	}
}

func TestReport_String_SomeIgnoredVulnerability(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  my-package@1.2.3 is affected by the following vulnerabilities:",
		"    GHSA-1: This is a vulnerability! (https://github.com/advisories/GHSA-1)",
		"",
		"  1 new vulnerability found in /path/to/my/lock (1 was ignored)",
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
					CompareAs: lockfile.BundlerEcosystem,
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
					Name:      "their-package",
					Version:   "4.5.6",
					Ecosystem: lockfile.BundlerEcosystem,
					CompareAs: lockfile.BundlerEcosystem,
				},
				Ignored: []database.OSV{
					{
						ID:      "GHSA-2",
						Summary: "This is another vulnerability!",
					},
				},
			},
		},
	}

	if actual := r.String(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}
