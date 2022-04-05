package reporter_test

import (
	"osv-detector/internal"
	"osv-detector/internal/database"
	"osv-detector/internal/lockfile"
	"osv-detector/internal/reporter"
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
