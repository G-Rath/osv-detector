package database

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ZipDB struct {
	vulnerabilities []OSV
	ArchiveURL      string
	Offline         bool
	UpdatedAt       string
}

// Cache stores the OSV database archive for re-use
type Cache struct {
	URL  string
	ETag string
	Date string
	Body []byte
}

var ErrOfflineDatabaseNotFound = errors.New("no offline version of the OSV database is available")

func (db *ZipDB) cachePath() string {
	hash := sha256.Sum256([]byte(db.ArchiveURL))
	fileName := fmt.Sprintf("osv-detector-%x-db.json", hash)

	return filepath.Join(os.TempDir(), fileName)
}

func (db *ZipDB) fetchZip() ([]byte, error) {
	var cache *Cache
	cachePath := db.cachePath()

	if cacheContent, err := os.ReadFile(cachePath); err == nil {
		err := json.Unmarshal(cacheContent, &cache)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to parse cache from %s: %v", cachePath, err)
		}
	}

	if db.Offline {
		if cache == nil {
			return nil, ErrOfflineDatabaseNotFound
		}

		db.UpdatedAt = cache.Date

		return cache.Body, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", db.ArchiveURL, nil)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve OSV database archive: %w", err)
	}

	if cache != nil {
		req.Header.Add("If-None-Match", cache.ETag)
		req.Header.Add("If-Modified-Since", cache.Date)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve OSV database archive: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		db.UpdatedAt = cache.Date

		return cache.Body, nil
	}

	var body []byte

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read OSV database archive from response: %w", err)
	}

	etag := resp.Header.Get("ETag")
	date := resp.Header.Get("Date")

	db.UpdatedAt = date

	if etag != "" || date != "" {
		cache = &Cache{ETag: etag, Date: date, Body: body, URL: db.ArchiveURL}
	}

	cacheContents, err := json.Marshal(cache)

	if err == nil {
		// nolint:gosec // being world readable is fine
		err = os.WriteFile(cachePath, cacheContents, 0644)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to write cache to %s: %v", cachePath, err)
		}
	}

	return body, nil
}

// load fetches a zip archive of the OSV database and loads known vulnerabilities
// from it (which are assumed to be in json files following the OSV spec).
//
// Internally, the archive is cached along with the date that it was fetched
// so that a new version of the archive is only downloaded if it has been
// modified, per HTTP caching standards.
func (db *ZipDB) load() error {
	db.vulnerabilities = []OSV{}

	body, err := db.fetchZip()

	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("could not read OSV database archive: %w", err)
	}

	// Read all the files from the zip archive
	for _, zipFile := range zipReader.File {
		if strings.HasPrefix(zipFile.Name, "advisory-database-main/advisories/unreviewed") {
			continue
		}

		if !strings.HasSuffix(zipFile.Name, ".json") {
			continue
		}

		file, err := zipFile.Open()
		if err != nil {
			return fmt.Errorf("could not open OSV database archive: %w", err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not read %s: %v", zipFile.Name, err)

			continue
		}

		var pa OSV
		if err := json.Unmarshal(content, &pa); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s is not a valid JSON file: %v", zipFile.Name, err)

			continue
		}

		db.vulnerabilities = append(db.vulnerabilities, pa)
	}

	return nil
}

func NewZippedDB(dbArchiveURL string, offline bool) (*ZipDB, error) {
	db := &ZipDB{Offline: offline, ArchiveURL: dbArchiveURL}
	if err := db.load(); err != nil {
		return nil, fmt.Errorf("unable to fetch OSV database: %w", err)
	}

	return db, nil
}
