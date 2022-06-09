package database_test

import (
	"net/url"
	"osv-detector/internal/database"
	"reflect"
	"testing"
)

func TestNewAPIDB(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://my-api.com")

	type args struct {
		baseURL   string
		batchSize int
		offline   bool
	}
	tests := []struct {
		name    string
		args    args
		want    *database.APIDB
		wantErr bool
	}{
		{
			name: "Offline is true",
			args: args{
				baseURL:   "https://my-api.com",
				batchSize: 100,
				offline:   true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Batch size is less than 1",
			args: args{
				baseURL:   "https://my-api.com",
				batchSize: 0,
				offline:   false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "URL is not valid",
			args: args{
				baseURL:   "not-a-url",
				batchSize: 100,
				offline:   false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Everything is valid",
			args: args{
				baseURL:   "https://my-api.com",
				batchSize: 100,
				offline:   false,
			},
			want:    &database.APIDB{BaseURL: u, BatchSize: 100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := database.NewAPIDB(tt.args.baseURL, tt.args.batchSize, tt.args.offline)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewAPIDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}
