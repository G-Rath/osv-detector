package database_test

import (
	"errors"
	"testing"

	"github.com/g-rath/osv-detector/pkg/database"
)

func TestNewDirDB(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{
		withDefaultAffected("OSV-1"),
		withDefaultAffected("OSV-2"),
		{
			ID: "GHSA-1234",
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "npm", Name: "request"}},
				{Package: database.Package{Ecosystem: "npm", Name: "@cypress/request"}},
			},
		},
	}

	db, err := database.NewDirDB(database.Config{URL: "file:/fixtures/db"}, false)

	if err != nil {
		t.Errorf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, osvs)
}

func TestNewDirDB_InvalidURI(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "file://\\"}, false)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}

func TestNewDirDB_NotFileProtocol(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "https://mysite.com/my.zip"}, false)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if !errors.Is(err, database.ErrDirPathWrongProtocol) {
		t.Errorf("NewDirDB() returned wrong error %v", err)
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}

func TestNewDirDB_DoesNotExist(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "file:/fixtures/nowhere"}, false)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}

func TestNewDirDB_WorkingDirectory(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{withDefaultAffected("OSV-1")}

	db, err := database.NewDirDB(database.Config{URL: "file:/fixtures/db", WorkingDirectory: "nested-1"}, false)

	if err != nil {
		t.Errorf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, osvs)
}
