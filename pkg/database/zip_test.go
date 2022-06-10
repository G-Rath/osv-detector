package database_test

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"osv-detector/pkg/database"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

type CleanUpZipServerFn = func()

func cachePath(url string) string {
	hash := sha256.Sum256([]byte(url))
	fileName := fmt.Sprintf("osv-detector-%x-db.json", hash)

	return filepath.Join(os.TempDir(), fileName)
}

func createZipServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, CleanUpZipServerFn) {
	t.Helper()

	ts := httptest.NewServer(handler)

	return ts, func() {
		ts.Close()

		_ = os.Remove(cachePath(ts.URL))
	}
}

func zipOSVs(t *testing.T, osvs map[string]database.OSV) []byte {
	t.Helper()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for filepath, osv := range osvs {
		data, err := json.Marshal(osv)
		if err != nil {
			t.Fatalf("could not marshal %v: %v", osv, err)
		}

		f, err := writer.Create(filepath)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

func TestNewZippedDB_OfflineWithoutCache(t *testing.T) {
	t.Parallel()

	ts, cleanup := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("a server request was made when running offline")
	})
	defer cleanup()

	_, err := database.NewZippedDB(ts.URL, true)

	if !errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		t.Errorf("expected \"%v\" error but got \"%v\"", database.ErrOfflineDatabaseNotFound, err)
	}
}

func TestNewZippedDB_BadZip(t *testing.T) {
	t.Parallel()

	ts, cleanup := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("this is not a zip"))
	})
	defer cleanup()

	_, err := database.NewZippedDB(ts.URL, false)

	if err == nil {
		t.Errorf("expected an error but did not get one")
	}
}

func TestNewZippedDB_UnsupportedProtocol(t *testing.T) {
	t.Parallel()

	_, err := database.NewZippedDB("file://hello-world", false)

	if err == nil {
		t.Errorf("expected an error but did not get one")
	}
}

func TestNewZippedDB_Online_Successful(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{
		{ID: "GHSA-1"},
		{ID: "GHSA-2"},
		{ID: "GHSA-3"},
		{ID: "GHSA-4"},
		{ID: "GHSA-5"},
	}

	ts, cleanup := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zipOSVs(t, map[string]database.OSV{
			"GHSA-1.json": {ID: "GHSA-1"},
			"GHSA-2.json": {ID: "GHSA-2"},
			"GHSA-3.json": {ID: "GHSA-3"},
			"GHSA-4.json": {ID: "GHSA-4"},
			"GHSA-5.json": {ID: "GHSA-5"},
		}))
	})
	defer cleanup()

	db, err := database.NewZippedDB(ts.URL, false)

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

func TestNewZippedDB_FileChecks(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{{ID: "GHSA-1234"}}

	ts, cleanup := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zipOSVs(t, map[string]database.OSV{
			"file.json": {ID: "GHSA-1234"},
			// only files with .json suffix should be loaded
			"file.yaml": {ID: "GHSA-5678"},
			// special case for the GH security database
			"advisory-database-main/advisories/unreviewed/file.json": {ID: "GHSA-4321"},
		}))
	})
	defer cleanup()

	db, err := database.NewZippedDB(ts.URL, false)

	if err != nil {
		t.Errorf("unexpected error \"%v\"", err)
	}

	if !reflect.DeepEqual(db.Vulnerabilities(false), osvs) {
		t.Errorf("db is missing some vulnerabilities: %v vs %v", db.Vulnerabilities(false), osvs)
	}
}
