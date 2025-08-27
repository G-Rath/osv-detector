package database_test

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/g-rath/osv-detector/pkg/database"
)

func dirPath(url string) string {
	hash := sha256.Sum256([]byte(url))
	fileName := fmt.Sprintf("osv-detector-%x-db", hash)

	return filepath.Join(os.TempDir(), fileName)
}

func dirWrite(t *testing.T, url string, lastModified string, osvs map[string]database.OSV) {
	t.Helper()

	err := os.Mkdir(dirPath(url), 0744)

	if err != nil {
		t.Fatalf("unexpected error creating directory: %v", err)
	}

	for _, osv := range osvs {
		f, err := os.OpenFile(dirPath(url)+"/"+osv.ID+".json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)

		if err != nil {
			t.Fatalf("unexpected error creating file: %v", err)
		}

		err = json.NewEncoder(f).Encode(osv)
		if err != nil {
			t.Fatalf("unexpected error encoding file: %v", err)
		}
	}

	err = os.WriteFile(dirPath(url)+"/last_checked", []byte(lastModified), 0600)

	if err != nil {
		t.Fatalf("error creating last_checked file: %v", err)
	}
}

func expectLastCheckedTime(t *testing.T, url string, expected time.Time) {
	t.Helper()

	expectedTime := expected.UTC().Format(time.RFC3339)
	checkedTime, err := os.ReadFile(dirPath(url) + "/last_checked")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(checkedTime) != expectedTime {
		t.Errorf("last checked time got = \"%s\", want = \"%s\"", string(checkedTime), expectedTime)
	}
}

func TestNewSmartDB_Offline_WithoutCache(t *testing.T) {
	t.Parallel()

	ts := createZipServer(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("a server request was made when running offline")
	})

	_, err := database.NewSmartDB(database.Config{URL: ts.URL}, true)

	if !errors.Is(err, database.ErrOfflineDatabaseNotFound) {
		t.Errorf("expected \"%v\" error but got \"%v\"", database.ErrOfflineDatabaseNotFound, err)
	}
}

func TestNewSmartDB_Offline_WithZipCache(t *testing.T) {
	t.Parallel()

	date := time.Now().UTC()
	osvs := []database.OSV{
		withDefaultAffected("GHSA-1"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withDefaultAffected("GHSA-4"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("a server request was made when running offline")
	})

	cacheWrite(t, database.Cache{
		URL:  ts.URL,
		ETag: "",
		Date: date.Format(http.TimeFormat),
		Body: zipOSVs(t, map[string]database.OSV{
			"GHSA-1.json": withDefaultAffected("GHSA-1"),
			"GHSA-2.json": withDefaultAffected("GHSA-2"),
			"GHSA-3.json": withDefaultAffected("GHSA-3"),
			"GHSA-4.json": withDefaultAffected("GHSA-4"),
			"GHSA-5.json": withDefaultAffected("GHSA-5"),
		}),
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, true)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if db.UpdatedAt != date.Format(http.TimeFormat) {
		t.Errorf("db.UpdatedAt got = \"%s\", want = \"%s\"", db.UpdatedAt, date)
	}

	expectDBToHaveOSVs(t, db, osvs)
	expectLastCheckedTime(t, ts.URL, date)
}

func TestNewSmartDB_Offline_WithDirCache(t *testing.T) {
	t.Parallel()

	date := time.Date(2022, time.June, 17, 22, 28, 13, 0, time.UTC)
	osvs := []database.OSV{
		withDefaultAffected("GHSA-1"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withDefaultAffected("GHSA-4"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("a server request was made when running offline")
	})

	cacheWrite(t, database.Cache{
		URL:  ts.URL,
		ETag: "",
		Date: date.Format(http.TimeFormat),
		Body: zipOSVs(t, map[string]database.OSV{
			"GHSA-1.json": withDefaultAffected("GHSA-1"),
			"GHSA-2.json": withDefaultAffected("GHSA-2"),
			"GHSA-3.json": withDefaultAffected("GHSA-3"),
			"GHSA-4.json": withDefaultAffected("GHSA-4"),
			"GHSA-5.json": withDefaultAffected("GHSA-5"),
		}),
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, true)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if db.UpdatedAt != date.Format(http.TimeFormat) {
		t.Errorf("db.UpdatedAt got = \"%s\", want = \"%s\"", db.UpdatedAt, date)
	}

	expectDBToHaveOSVs(t, db, osvs)

	// when offline, the "last checked" time should not be changed
	expectLastCheckedTime(t, ts.URL, date)
}

func TestNewSmartDB_BadZip(t *testing.T) {
	t.Parallel()

	ts := createZipServer(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("this is not a zip"))
	})

	_, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err == nil {
		t.Errorf("expected an error but did not get one")
	}
}

func TestNewSmartDB_UnsupportedProtocol(t *testing.T) {
	t.Parallel()

	_, err := database.NewSmartDB(database.Config{URL: "file://hello-world"}, false)

	if err == nil {
		t.Errorf("expected an error but did not get one")
	}
}

func TestNewSmartDB_Online_WithoutCache(t *testing.T) {
	t.Parallel()

	osvs := []database.OSV{
		withDefaultAffected("GHSA-1"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withDefaultAffected("GHSA-4"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/modified_id.csv" {
			t.Errorf("unexpected request for modified_id.csv")

			_, _ = w.Write([]byte(""))

			return
		}

		_, _ = w.Write(zipOSVs(t, map[string]database.OSV{
			"GHSA-1.json": withDefaultAffected("GHSA-1"),
			"GHSA-2.json": withDefaultAffected("GHSA-2"),
			"GHSA-3.json": withDefaultAffected("GHSA-3"),
			"GHSA-4.json": withDefaultAffected("GHSA-4"),
			"GHSA-5.json": withDefaultAffected("GHSA-5"),
		}))
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	expectDBToHaveOSVs(t, db, osvs)
}

func TestNewSmartDB_Online_WithoutCache_NotFound(t *testing.T) {
	t.Parallel()

	ts := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/modified_id.csv" {
			t.Errorf("unexpected request for modified_id.csv")

			_, _ = w.Write([]byte(""))

			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(zipOSVs(t, map[string]database.OSV{}))
	})

	_, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err == nil {
		t.Errorf("expected an error but did not get one")
	} else if !errors.Is(err, database.ErrUnexpectedStatusCode) {
		t.Errorf("expected %v error but got %v", database.ErrUnexpectedStatusCode, err)
	}
}

func toJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestNewSmartDB_Online_WithExistingDirDB_UpToDate(t *testing.T) {
	t.Parallel()

	date := time.Now().Add(-14 * time.Hour).UTC()
	osvs := []database.OSV{
		withDefaultAffected("GHSA-1"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withDefaultAffected("GHSA-4"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/modified_id.csv" {
			_, _ = w.Write([]byte(strings.Join([]string{
				date.Add(-02*time.Hour).Format(time.RFC3339) + ",GHSA-4",
				date.Add(-38*time.Minute).Format(time.RFC3339) + ",GHSA-2",
				date.Add(-06*time.Hour).Format(time.RFC3339) + ",GHSA-1",
				date.Add(-03*time.Minute).Format(time.RFC3339) + ",GHSA-3",
				date.Add(-14*time.Hour).Format(time.RFC3339) + ",GHSA-5",
			}, "\n")))

			return
		}

		t.Errorf("unexpected request for %s", r.URL.Path)

		_, _ = w.Write([]byte("{}"))
	})

	dirWrite(t, ts.URL, date.Format(time.RFC3339), map[string]database.OSV{
		"GHSA-1.json": withDefaultAffected("GHSA-1"),
		"GHSA-2.json": withDefaultAffected("GHSA-2"),
		"GHSA-3.json": withDefaultAffected("GHSA-3"),
		"GHSA-4.json": withDefaultAffected("GHSA-4"),
		"GHSA-5.json": withDefaultAffected("GHSA-5"),
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if now := time.Now().UTC().Format(http.TimeFormat); db.UpdatedAt != now {
		t.Errorf("db.UpdatedAt got = \"%s\", want = \"%s\"", db.UpdatedAt, now)
	}

	expectDBToHaveOSVs(t, db, osvs)

	// when online, the "last checked" time should be now
	expectLastCheckedTime(t, ts.URL, time.Now().UTC())
}

func TestNewSmartDB_Online_WithExistingDirDB_NotModified(t *testing.T) {
	t.Parallel()

	date := time.Now().Add(-14 * time.Hour).UTC()
	osvs := []database.OSV{
		withDefaultAffected("GHSA-1"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withDefaultAffected("GHSA-4"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/modified_id.csv" {
			if dateHeader := r.Header.Get("If-Modified-Since"); dateHeader != date.Format(http.TimeFormat) {
				t.Errorf("incorrect Date header: got = \"%s\", want = \"%s\"", dateHeader, date)
			}

			w.WriteHeader(http.StatusNotModified)

			return
		}

		t.Errorf("unexpected request for %s", r.URL.Path)

		_, _ = w.Write([]byte("{}"))
	})

	dirWrite(t, ts.URL, date.Format(time.RFC3339), map[string]database.OSV{
		"GHSA-1.json": withDefaultAffected("GHSA-1"),
		"GHSA-2.json": withDefaultAffected("GHSA-2"),
		"GHSA-3.json": withDefaultAffected("GHSA-3"),
		"GHSA-4.json": withDefaultAffected("GHSA-4"),
		"GHSA-5.json": withDefaultAffected("GHSA-5"),
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if now := time.Now().UTC().Format(http.TimeFormat); db.UpdatedAt != now {
		t.Errorf("db.UpdatedAt got = \"%s\", want = \"%s\"", db.UpdatedAt, now)
	}

	expectDBToHaveOSVs(t, db, osvs)

	// when online, the "last checked" time should be now
	expectLastCheckedTime(t, ts.URL, time.Now().UTC())
}

func TestNewSmartDB_Online_WithExistingDirDB_Outdated(t *testing.T) {
	t.Parallel()

	date := time.Now().Add(-14 * time.Hour).UTC()
	osvs := []database.OSV{
		withSummary("GHSA-1", "this summary has changed"),
		withDefaultAffected("GHSA-2"),
		withDefaultAffected("GHSA-3"),
		withSummary("GHSA-4", "and so has this one!"),
		withDefaultAffected("GHSA-5"),
	}

	ts := createZipServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/modified_id.csv":
			_, _ = w.Write([]byte(strings.Join([]string{
				date.Add(+02*time.Hour).Format(time.RFC3339) + ",GHSA-4",
				date.Add(+38*time.Minute).Format(time.RFC3339) + ",GHSA-2",
				date.Add(+06*time.Hour).Format(time.RFC3339) + ",GHSA-1",
				date.Add(-03*time.Minute).Format(time.RFC3339) + ",GHSA-3",
				date.Add(-14*time.Hour).Format(time.RFC3339) + ",GHSA-5",
			}, "\n")))

			return
		case "/GHSA-1.json":
			_, _ = w.Write(toJSON(t, withSummary("GHSA-1", "this summary has changed")))
		case "/GHSA-2.json":
			_, _ = w.Write(toJSON(t, withDefaultAffected("GHSA-2")))
		case "/GHSA-3.json":
			t.Errorf("unexpected request for GHSA-3.json")
			_, _ = w.Write(toJSON(t, withDefaultAffected("GHSA-3")))
		case "/GHSA-4.json":
			_, _ = w.Write(toJSON(t, withSummary("GHSA-4", "and so has this one!")))
		case "/GHSA-5.json":
			t.Errorf("unexpected request for GHSA-5.json")
			_, _ = w.Write(toJSON(t, withDefaultAffected("GHSA-5")))
		}
	})

	dirWrite(t, ts.URL, date.Format(time.RFC3339), map[string]database.OSV{
		"GHSA-1.json": withDefaultAffected("GHSA-1"),
		"GHSA-2.json": withDefaultAffected("GHSA-2"),
		"GHSA-3.json": withDefaultAffected("GHSA-3"),
		"GHSA-4.json": withDefaultAffected("GHSA-4"),
		"GHSA-5.json": withDefaultAffected("GHSA-5"),
	})

	db, err := database.NewSmartDB(database.Config{URL: ts.URL}, false)

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if now := time.Now().UTC().Format(http.TimeFormat); db.UpdatedAt != now {
		t.Errorf("db.UpdatedAt got = \"%s\", want = \"%s\"", db.UpdatedAt, now)
	}

	expectDBToHaveOSVs(t, db, osvs)

	// when online, the "last checked" time should be now
	expectLastCheckedTime(t, ts.URL, time.Now().UTC())
}
