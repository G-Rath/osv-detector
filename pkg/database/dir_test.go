package database_test

import (
	"errors"
	"osv-detector/pkg/database"
	"reflect"
	"sort"
	"testing"
)

func TestNewDirDB(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{{ID: "OSV-1"}, {ID: "OSV-2"}, {ID: "GHSA-1234"}}

	db, err := database.NewDirDB(database.Config{URL: "file:/fixtures/db"}, false)

	if err != nil {
		t.Errorf("unexpected error \"%v\"", err)
	}

	vulns := db.Vulnerabilities(false)

	sort.Slice(vulns, func(i, j int) bool {
		return vulns[i].ID < vulns[j].ID
	})
	sort.Slice(osvs, func(i, j int) bool {
		return osvs[i].ID < osvs[j].ID
	})

	if !reflect.DeepEqual(vulns, osvs) {
		t.Errorf("db is missing some vulnerabilities: %v vs %v", vulns, osvs)
	}
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

	osvs := []database.OSV{{ID: "OSV-1"}}

	db, err := database.NewDirDB(database.Config{URL: "file:/fixtures/db", WorkingDirectory: "nested-1"}, false)

	if err != nil {
		t.Errorf("unexpected error \"%v\"", err)
	}

	if !reflect.DeepEqual(db.Vulnerabilities(false), osvs) {
		t.Errorf("db is missing some vulnerabilities: %v vs %v", db.Vulnerabilities(false), osvs)
	}
}
