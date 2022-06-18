package database_test

import (
	"osv-detector/pkg/database"
	"testing"
)

func TestLoad(t *testing.T) {
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
	db, err := database.Load(database.Config{Type: "file"}, false, 100)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}
