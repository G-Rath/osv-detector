package database_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/g-rath/osv-detector/pkg/database"
)

func TestNewAPIDB(t *testing.T) {
	t.Parallel()

	type args struct {
		baseURL   string
		batchSize int
		offline   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Offline is true",
			args: args{
				baseURL:   "https://my-api.com",
				batchSize: 100,
				offline:   true,
			},
		},
		{
			name: "Batch size is less than 1",
			args: args{
				baseURL:   "https://my-api.com",
				batchSize: 0,
				offline:   false,
			},
		},
		{
			name: "URL is not valid",
			args: args{
				baseURL:   "not-a-url",
				batchSize: 100,
				offline:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := database.NewAPIDB(database.Config{URL: tt.args.baseURL}, tt.args.offline, tt.args.batchSize)
			if err == nil {
				t.Errorf("NewAPIDB() did not error as expected")
			}
			if got != nil {
				t.Errorf("NewAPIDB() returned a db even though there was an error")
			}
		})
	}
}

func TestNewAPIDB_Valid(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("https://my-api.com")

	config := database.Config{URL: "https://my-api.com", Name: "my-api"}

	db, err := database.NewAPIDB(
		config,
		false,
		100,
	)

	if err != nil {
		t.Errorf("NewAPIDB() unexpected error \"%v\"", err)
	}

	if db == nil {
		t.Fatalf("NewAPIDB() db unexpectedly nil")

		// this is required currently to make the staticcheck linter
		return
	}

	if !reflect.DeepEqual(db.BaseURL, u) {
		t.Errorf("NewAPIDB() db has incorrect url (%s)", u)
	}

	if db.BatchSize != 100 {
		t.Errorf("NewAPIDB() db has incorrect batch size (%d)", db.BatchSize)
	}

	if db.Identifier() != config.Identifier() {
		t.Errorf("NewAPIDB() db identifier got = %s, want %s", db.Identifier(), config.Identifier())
	}

	if db.Name() != config.Name {
		t.Errorf("NewAPIDB() db name got = %s, want %s", db.Name(), config.Name)
	}
}
