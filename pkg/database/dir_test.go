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
			ID: "OSV-3",
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "PyPi", Name: "mine2"}, Versions: database.Versions{}},
			},
		},
		{
			ID: "GHSA-1234",
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "npm", Name: "request"}},
				{Package: database.Package{Ecosystem: "npm", Name: "@cypress/request"}},
			},
		},
	}

	db, err := database.NewDirDB(database.Config{URL: "file:/testdata/db"}, false, nil)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, osvs)
}

func TestNewDirDB_InvalidURI(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "file://\\"}, false, nil)

	if err == nil {
		t.Fatalf("NewDirDB() did not return expected error")
	}

	if db != nil {
		t.Errorf("NewDirDB() returned a db even though it errored")
	}
}

func TestNewDirDB_NotFileProtocol(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "https://mysite.com/my.zip"}, false, nil)

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

	db, err := database.NewDirDB(database.Config{URL: "file:/testdata/nowhere"}, false, nil)

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

	db, err := database.NewDirDB(database.Config{URL: "file:/testdata/db", WorkingDirectory: "nested-1"}, false, nil)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, osvs)
}

func TestNewDirDB_WithSpecificPackages(t *testing.T) {
	t.Parallel()

	db, err := database.NewDirDB(database.Config{URL: "file:/testdata/db"}, false, []string{"mine", "request"})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, []database.OSV{
		withDefaultAffected("OSV-1"),
		withDefaultAffected("OSV-2"),
		{
			ID: "GHSA-1234",
			Affected: []database.Affected{
				{Package: database.Package{Ecosystem: "npm", Name: "request"}},
				{Package: database.Package{Ecosystem: "npm", Name: "@cypress/request"}},
			},
		},
	})
}
