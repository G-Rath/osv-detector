package database

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SmartDB struct {
	// memDB
	DirDB

	name             string
	identifier       string
	ArchiveURL       string
	WorkingDirectory string
	Offline          bool
	UpdatedAt        string

	cacheDirectory string
}

func (db *SmartDB) Name() string       { return db.name }
func (db *SmartDB) Identifier() string { return db.identifier }

func (db *SmartDB) cachePath() string {
	hash := sha256.Sum256([]byte(db.ArchiveURL))
	fileName := fmt.Sprintf("osv-detector-%x-db", hash)

	return filepath.Join(db.cacheDirectory, fileName)
}

func (db *SmartDB) cacheFile(name string, content []byte) error {
	//nolint:gosec // being world readable is fine
	return os.WriteFile(filepath.Join(db.cachePath(), name), content, 0644)
}

func (db *SmartDB) loadLastChecked() (time.Time, error) {
	b, err := os.ReadFile(filepath.Join(db.cachePath(), "last_checked"))

	if err != nil {
		return time.Time{}, err
	}

	tim, err := time.Parse(time.RFC3339, string(b))

	if err != nil {
		return time.Time{}, err
	}

	return tim, nil
}

func (db *SmartDB) updateLastChecked(lastChecked time.Time) error {
	db.UpdatedAt = lastChecked.Format(http.TimeFormat)

	return db.cacheFile("last_checked", []byte(lastChecked.Format(time.RFC3339)))
}

func (db *SmartDB) writeZipFile(zipFile *zip.File) error {
	dst, err := os.OpenFile(filepath.Join(db.cachePath(), zipFile.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}

	defer dst.Close()

	z, err := zipFile.Open()
	if err != nil {
		return err
	}

	defer z.Close()

	_, err = io.Copy(dst, z)

	return err
}

func (db *SmartDB) populateFromZip() (*time.Time, error) {
	err := os.MkdirAll(db.cachePath(), 0744)

	if err != nil {
		return nil, err
	}

	zdb := &ZipDB{
		name:             db.name,
		ArchiveURL:       db.ArchiveURL,
		WorkingDirectory: db.WorkingDirectory,
		cacheDirectory:   db.cacheDirectory,
		Offline:          db.Offline,
	}

	body, err := zdb.fetchZip()

	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("could not read OSV database archive: %w", err)
	}

	// Read each file from the archive and write it to the db directory
	for _, zipFile := range zipReader.File {
		if !strings.HasPrefix(zipFile.Name, db.WorkingDirectory) {
			continue
		}

		if !strings.HasSuffix(zipFile.Name, ".json") {
			continue
		}

		err = db.writeZipFile(zipFile)

		if err != nil {
			return nil, err
		}
	}

	tim, err := time.Parse(http.TimeFormat, zdb.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &tim, nil
}

type modifiedIDRow struct {
	id       string
	modified time.Time
}

func parseModifiedIDRow(columns []string) (*modifiedIDRow, error) {
	modified, err := time.Parse(time.RFC3339, columns[0])

	if err != nil {
		return nil, err
	}

	return &modifiedIDRow{id: columns[1], modified: modified}, nil
}

func (db *SmartDB) fetchModifiedIDs(since time.Time) ([]string, error) {
	url := strings.TrimSuffix(db.ArchiveURL, "/all.zip") + "/modified_id.csv"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("If-Modified-Since", since.Format(http.TimeFormat))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w (%s)", ErrUnexpectedStatusCode, resp.Status)
	}

	i := 0
	r := csv.NewReader(resp.Body)

	var ids []string

	for {
		i++
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		row, err := parseModifiedIDRow(record)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i, err)
		}

		// the modified ids are sorted in reverse chronological order so once we hit
		// a row that was modified before our "since" time, we can stop completely
		if row.modified.Before(since) {
			break
		}

		ids = append(ids, row.id)
	}

	return ids, nil
}

func (db *SmartDB) updateAdvisory(id string) error {
	url := fmt.Sprintf("%s/%s.json", strings.TrimSuffix(db.ArchiveURL, "/all.zip"), id)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w (%s)", ErrUnexpectedStatusCode, resp.Status)
	}

	var body []byte

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = db.cacheFile(id+".json", body)

	return err
}

func (db *SmartDB) downloadModifiedAdvisories(ids []string) error {
	conLimit := 200

	if len(ids) == 0 {
		return nil
	}

	// buffered channel which controls the number of concurrent operations
	semaphoreChan := make(chan struct{}, conLimit)
	resultsChan := make(chan *result)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, id := range ids {
		go func(i int, id string) {
			// read from the buffered semaphore channel, which will block if we're
			// already got as many goroutines as our concurrency limit allows
			//
			// when one of those routines finish they'll read from this channel,
			// freeing up a slot to unblock this send
			semaphoreChan <- struct{}{}

			// use an empty OSV as we're reusing the result struct
			result := &result{i, OSV{}, db.updateAdvisory(id)}

			resultsChan <- result

			// read from the buffered semaphore to free up a slot to allow
			// another goroutine to start, since this one is wrapping up
			<-semaphoreChan
		}(i, id)
	}

	var errs []error

	for {
		result := <-resultsChan
		errs = append(errs, result.err)

		if len(errs) == len(ids) {
			break
		}
	}

	return errors.Join(errs...)
}

func (db *SmartDB) updateModifiedAdvisories(since time.Time) error {
	modifiedIDs, err := db.fetchModifiedIDs(since)

	if err != nil {
		return err
	}

	return db.downloadModifiedAdvisories(modifiedIDs)
}

func (db *SmartDB) populate() (*time.Time, error) {
	lastChecked, err := db.loadLastChecked()

	// if there's an error, assumingly because the database does not already exist
	// then extract it from a zip file, and use the zips updated date as the timestamp
	if err != nil {
		return db.populateFromZip()
	}

	// if we're offline, then we can only work with the database on-hand, and don't
	// want to change the last checked time
	if db.Offline {
		return &lastChecked, nil
	}

	// otherwise, update all the advisories that have changed since our last check
	err = db.updateModifiedAdvisories(lastChecked)
	if err != nil {
		return nil, err
	}

	lastChecked = time.Now().UTC()

	return &lastChecked, nil
}

// load fetches a zip archive of the OSV database and loads known vulnerabilities
// from it (which are assumed to be in json files following the OSV spec).
//
// Internally, the archive is cached along with the date that it was fetched
// so that a new version of the archive is only downloaded if it has been
// modified, per HTTP caching standards.
func (db *SmartDB) load() error {
	lastChecked, err := db.populate()

	if err != nil {
		return err
	}

	if err = db.updateLastChecked(*lastChecked); err != nil {
		return err
	}

	db.DirDB = DirDB{
		name:             db.name,
		LocalPath:        "file:///" + db.cachePath(),
		WorkingDirectory: "",
		Offline:          db.Offline,
	}

	return db.DirDB.load()
}

func NewSmartDB(config Config, offline bool) (*SmartDB, error) {
	if config.CacheDirectory == "" {
		d, err := setupCacheDirectory()

		if err != nil {
			return nil, err
		}

		config.CacheDirectory = d
	}

	db := &SmartDB{
		name:             config.Name,
		identifier:       config.Identifier(),
		ArchiveURL:       config.URL,
		WorkingDirectory: config.WorkingDirectory,
		cacheDirectory:   config.CacheDirectory,
		Offline:          offline,
	}
	if err := db.load(); err != nil {
		return nil, fmt.Errorf("unable to fetch OSV database: %w", err)
	}

	return db, nil
}
