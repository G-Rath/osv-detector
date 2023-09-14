package database_test

import (
	"errors"
	"testing"

	"github.com/g-rath/osv-detector/pkg/database"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	types := []string{
		"zip",
		"dir",
		"api",
	}

	for _, typ := range types {
		_, err := database.Load(database.Config{Type: typ}, false, 100)

		if err == nil {
			t.Fatalf("NewDirDB() did not return expected error")
		}
	}
}

func TestLoad_BadType(t *testing.T) {
	t.Parallel()

	db, err := database.Load(database.Config{Type: "file"}, false, 100)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if !errors.Is(err, database.ErrUnsupportedDatabaseType) {
		t.Errorf("NewDirDB() returned wrong error %v", err)
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}
