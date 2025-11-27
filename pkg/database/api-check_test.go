package database_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/lockfile"
)

type objectsWithIDs = []database.ObjectWithID
type apiPackage struct {
	Name      string             `json:"name"`
	Ecosystem internal.Ecosystem `json:"ecosystem"`
}
type apiQuery struct {
	Commit  string     `json:"commit,omitempty"`
	Version string     `json:"version,omitempty"`
	Package apiPackage `json:"package"`
}

func jsonMarshalQueryBatchResponse(t *testing.T, vulns []objectsWithIDs) []byte {
	t.Helper()

	type vulnsType struct {
		Vulns objectsWithIDs `json:"vulns"`
	}
	type results struct {
		Results []vulnsType `json:"results"`
	}

	var payload results

	for _, actualVulns := range vulns {
		payload.Results = append(payload.Results, vulnsType{actualVulns})
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		t.Fatalf("could not marshal test server response: %s", err)
	}

	return jsonData
}

func expectRequestPayload(t *testing.T, r *http.Request, queries []apiQuery) {
	t.Helper()

	if r.Method != http.MethodPost {
		t.Fatalf("api query was not a POST request")
	}

	var payload struct {
		Queries []apiQuery `json:"queries"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		t.Fatalf("could not decode api request payload: %s", err)
	}

	if len(queries) != len(payload.Queries) {
		t.Errorf("queries count mismatch: %d vs %d", len(queries), len(payload.Queries))
	}

	if !reflect.DeepEqual(payload.Queries, queries) {
		t.Errorf("queries mismatch: got %v, want %v", payload.Queries, queries)
	}
}

func expectVulnerability(t *testing.T, vuln database.OSV, id string, summary string) {
	t.Helper()

	if vuln.ID != id {
		t.Errorf("expected vulnerability id to be \"%s\" but got \"%s\"", id, vuln.ID)
	}

	if vuln.Summary != summary {
		t.Errorf("expected vulnerability to have summary \"%s\" but instead got \"%s\"", summary, vuln.Summary)
	}
}

func TestAPIDB_Check_NoPackages(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Errorf("an API request was made even though there are no packages to check")
	}))
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) > 0 {
		t.Fatalf("did not expect to get any packages, but got %d", len(vulns))
	}
}

func TestAPIDB_Check_UnsupportedProtocol(t *testing.T) {
	t.Parallel()

	db, err := database.NewAPIDB(database.Config{URL: "file://hello-world"}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "abc123", Ecosystem: "npm", CompareAs: "npm"},
	})

	if err == nil {
		t.Errorf("expected an error but did not get one")
	}

	if len(vulns) > 0 {
		t.Fatalf("did not expect to get any packages, but got %d", len(vulns))
	}
}

func TestAPIDB_Check_NotOK(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		_, _ = w.Write([]byte("<html></html>"))
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "", Ecosystem: "npm", CompareAs: "npm"},
	})

	if err == nil {
		t.Errorf("expected err but did not get one")
	}

	if len(vulns) > 0 {
		t.Fatalf("did not expect to get any packages, but got %d", len(vulns))
	}
}

func TestAPIDB_Check_InvalidBody(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		http.Error(w, "oh noes!", http.StatusForbidden)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "", Ecosystem: "npm", CompareAs: "npm"},
	})

	if err == nil {
		t.Errorf("expected err but did not get one")
	}

	if len(vulns) > 0 {
		t.Fatalf("did not expect to get any packages, but got %d", len(vulns))
	}
}

func TestAPIDB_Check_UnbalancedResponse(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
			{
				Version: "1.2.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			}},
		)

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{}})

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 2)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Ecosystem: "npm"},
		{Name: "my-package", Version: "1.2.0", Ecosystem: "npm"},
	})

	if !errors.Is(err, database.ErrAPIResultsCountMismatch) {
		t.Errorf("expected \"%v\" error but got \"%v\"", database.ErrAPIResultsCountMismatch, err)
	}

	if len(vulns) != 0 {
		t.Fatalf("expected to get 0 packages but got %d", len(vulns))
	}
}

func TestAPIDB_Check_FetchSuccessful(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{{"GHSA-1234"}}})

		_, _ = w.Write(jsonData)
	})

	mux.HandleFunc("/vulns/GHSA-1234", func(w http.ResponseWriter, _ *http.Request) {
		jsonData, err := json.Marshal(database.OSV{ID: "GHSA-1234", Summary: "my vulnerability"})

		if err != nil {
			t.Fatalf("could not marshal test server response: %s", err)
		}

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "", Ecosystem: "npm", CompareAs: "npm"},
	})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 1 {
		t.Errorf("expected to get 1 vulnerability but got %d", len(vulns[0]))
	}

	expectVulnerability(t, vulns[0][0], "GHSA-1234", "my vulnerability")
}

func TestAPIDB_Check_FetchFails(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{{"GHSA-1234"}, {"GHSA-5678"}}})

		_, _ = w.Write(jsonData)
	})

	// this response is not a 200 OK
	mux.HandleFunc("/vulns/GHSA-1234", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "oh noes!", http.StatusForbidden)
	})

	// this response is not valid json
	mux.HandleFunc("/vulns/GHSA-5678", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html></html>"))
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "", Ecosystem: "npm", CompareAs: "npm"},
	})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 2 {
		t.Errorf("expected to get 2 vulnerabilities but got %d", len(vulns[0]))
	}

	expectVulnerability(t, vulns[0][0], "GHSA-1234", "")
	expectVulnerability(t, vulns[0][1], "GHSA-5678", "")
}

func TestAPIDB_Check_FetchMixed(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{{"GHSA-1234"}, {"GHSA-5678"}}})

		_, _ = w.Write(jsonData)
	})

	mux.HandleFunc("/vulns/GHSA-1234", func(w http.ResponseWriter, _ *http.Request) {
		jsonData, err := json.Marshal(database.OSV{ID: "GHSA-1234", Summary: "my vulnerability"})

		if err != nil {
			t.Fatalf("could not marshal test server response: %s", err)
		}

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "", Ecosystem: "npm"},
	})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 2 {
		t.Errorf("expected to get 2 vulnerabilities but got %d", len(vulns[0]))
	}

	expectVulnerability(t, vulns[0][0], "GHSA-1234", "my vulnerability")
	expectVulnerability(t, vulns[0][1], "GHSA-5678", "")
}

func TestAPIDB_Check_WithCommit(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{
			{
				Version: "1.0.0",
				Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
			},
		})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{}})

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Commit: "abc123", Ecosystem: "npm"},
	})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 0 {
		t.Fatalf("expected to get 0 vulnerabilities but got %d", len(vulns[0]))
	}
}

func TestAPIDB_Check_WithCommitOnly(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{{Commit: "abc123"}})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{}})

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{{Commit: "abc123"}})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 0 {
		t.Fatalf("expected to get 0 vulnerabilities but got %d", len(vulns[0]))
	}
}

func TestAPIDB_Check_WithCommitAndSomeFields(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		expectRequestPayload(t, r, []apiQuery{{Commit: "abc123"}})

		jsonData := jsonMarshalQueryBatchResponse(t, []objectsWithIDs{{}})

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 1)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{{Commit: "abc123", Ecosystem: "npm"}})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if len(vulns) != 1 {
		t.Fatalf("expected to get 1 package but got %d", len(vulns))
	}

	if len(vulns[0]) != 0 {
		t.Fatalf("expected to get 0 vulnerabilities but got %d", len(vulns[0]))
	}
}

func TestAPIDB_Check_Batches(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	requestCount := 0

	mux.HandleFunc("/querybatch", func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount > 2 {
			t.Errorf("unexpected number of requests (%d)", requestCount)
		}

		var expectedPayload []apiQuery
		var batchResponse []objectsWithIDs

		// strictly speaking not the best of checks, but it should be good enough
		if r.ContentLength > 100 {
			expectedPayload = []apiQuery{
				{
					Version: "1.0.0",
					Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
				},
				{
					Version: "1.2.0",
					Package: apiPackage{Name: "my-package", Ecosystem: lockfile.NpmEcosystem},
				},
			}
			batchResponse = []objectsWithIDs{{}, {}}
		} else if r.ContentLength > 50 {
			expectedPayload = []apiQuery{
				{
					Version: "2.3.1",
					Package: apiPackage{Name: "their-package", Ecosystem: lockfile.NpmEcosystem},
				},
			}
			batchResponse = []objectsWithIDs{{}}
		}

		expectRequestPayload(t, r, expectedPayload)
		jsonData := jsonMarshalQueryBatchResponse(t, batchResponse)

		_, _ = w.Write(jsonData)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	db, err := database.NewAPIDB(database.Config{URL: ts.URL}, false, 2)

	if err != nil {
		t.Fatalf("Check() unexpected error \"%v\"", err)
	}

	vulns, err := db.Check([]internal.PackageDetails{
		{Name: "my-package", Version: "1.0.0", Ecosystem: "npm"},
		{Name: "my-package", Version: "1.2.0", Ecosystem: "npm"},
		{Name: "their-package", Version: "2.3.1", Ecosystem: "npm"},
	})

	if err != nil {
		t.Fatalf("unexpected error \"%v\"", err)
	}

	if requestCount != 2 {
		t.Errorf("expected there to be 2 requests but instead there were %d", requestCount)
	}

	if len(vulns) != 3 {
		t.Fatalf("expected to get 3 packages but got %d", len(vulns))
	}

	for _, vuln := range vulns {
		if len(vuln) != 0 {
			t.Errorf("expected to get 0 vulnerabilities but got %d", len(vuln))
		}
	}
}
