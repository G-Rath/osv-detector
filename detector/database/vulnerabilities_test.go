package database_test

import (
	"osv-detector/detector/database"
	"testing"
)

func TestVulnerabilities_Includes(t *testing.T) {
	t.Parallel()

	type args struct {
		osv database.OSV
	}
	tests := []struct {
		name string
		vs   database.Vulnerabilities
		args args
		want bool
	}{
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-2",
					Aliases: []string{},
				},
			},
			want: false,
		},
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{},
				},
			},
			want: true,
		},
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{"GHSA-2"},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-2",
					Aliases: []string{},
				},
			},
			want: true,
		},
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-2",
					Aliases: []string{"GHSA-1"},
				},
			},
			want: true,
		},
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{"CVE-1"},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-2",
					Aliases: []string{"CVE-1"},
				},
			},
			want: true,
		},
		{
			name: "",
			vs: database.Vulnerabilities{
				database.OSV{
					ID:      "GHSA-1",
					Aliases: []string{"CVE-2"},
				},
			},
			args: args{
				osv: database.OSV{
					ID:      "GHSA-2",
					Aliases: []string{"CVE-2"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.vs.Includes(tt.args.osv); got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}
