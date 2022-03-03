package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Cache stores the GitHub response to save bandwidth
type Cache struct {
	ETag string
	Date string
	Body []byte
}

func (db *OSVDatabase) fetchCache() (*Cache, error) {
	var cache *Cache
	cachePath := filepath.Join(os.TempDir(), "osv-detector-db.json")
	if cacheContent, err := ioutil.ReadFile(cachePath); err == nil {
		err := json.Unmarshal(cacheContent, &cache)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to parse cache from %s: %v", cachePath, err)
		}
	}

	if db.Offline && cache == nil {
		return nil, errors.New("--offline can only be used when a local version of the OSV database is available")
	}

	if !db.Offline {
		req, err := http.NewRequest("GET", db.ArchiveURL, nil)

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

		var body []byte

		if resp.StatusCode == http.StatusNotModified {
			return cache, nil
		}

		body, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("could not read OSV database archive from response: %w", err)
		}

		etag := resp.Header.Get("ETag")
		date := resp.Header.Get("Date")

		if etag != "" || date != "" {
			cache = &Cache{ETag: etag, Date: date, Body: body}
		}

		cacheContents, err := json.Marshal(cache)

		if err == nil {
			err = ioutil.WriteFile(cachePath, cacheContents, 0644)

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to write cache to %s: %v", cachePath, err)
			}
		}
	}

	return cache, nil
}
