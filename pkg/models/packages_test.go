package models_test

import (
	"github.com/g-rath/osv-detector/pkg/models"
	"reflect"
	"testing"
)

func TestPackages_Ecosystems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ps   models.Packages
		want []models.Ecosystem
	}{
		{name: "", ps: models.Packages{}, want: []models.Ecosystem{}},
		{
			name: "",
			ps: models.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: models.EcosystemCratesIO,
				},
			},
			want: []models.Ecosystem{
				models.EcosystemCratesIO,
			},
		},
		{
			name: "",
			ps: models.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: models.EcosystemCratesIO,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: models.EcosystemCratesIO,
				},
			},
			want: []models.Ecosystem{
				models.EcosystemCratesIO,
			},
		},
		{
			name: "",
			ps: models.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: models.EcosystemCratesIO,
				},
				{
					Name:      "@typescript-eslint/types",
					Version:   "5.13.0",
					Ecosystem: models.EcosystemNPM,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: models.EcosystemCratesIO,
				},
				{
					Name:      "sentry/sdk",
					Version:   "2.0.4",
					Ecosystem: models.EcosystemPackagist,
				},
			},
			want: []models.Ecosystem{
				models.EcosystemPackagist,
				models.EcosystemCratesIO,
				models.EcosystemNPM,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.ps.Ecosystems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ecosystems() = %v, want %v", got, tt.want)
			}
		})
	}
}
