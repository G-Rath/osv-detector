package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type DirDB struct {
	memDB

	name             string
	identifier       string
	LocalPath        string
	WorkingDirectory string
	Offline          bool
}

func (db *DirDB) Name() string       { return db.name }
func (db *DirDB) Identifier() string { return db.identifier }

var ErrDirPathWrongProtocol = errors.New("directory path must start with \"file:\" protocol")

// load walks the filesystem starting with the working directory within the local path,
// loading all OSVs found along the way.
func (db *DirDB) load(pkgNames []string) error {
	db.vulnerabilities = make(map[string][]OSV)

	if !strings.HasPrefix(db.LocalPath, "file:") {
		return ErrDirPathWrongProtocol
	}

	u, err := url.ParseRequestURI(db.LocalPath)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// since this will always be an absolute _url_ we need to remove the
	// absolute slash that will always be present
	p := filepath.Join(strings.TrimPrefix(u.Path, "/"), db.WorkingDirectory)
	errored := false

	err = filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}

		if err != nil {
			errored = true
			_, _ = fmt.Fprintf(os.Stderr, "\n    %v", err)

			return nil
		}

		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			errored = true
			_, _ = fmt.Fprintf(os.Stderr, "\n%v", err)

			return nil
		}

		var pa OSV
		if err := json.Unmarshal(content, &pa); err != nil {
			errored = true
			_, _ = fmt.Fprintf(os.Stderr, "%s is not a valid JSON file: %v\n", info.Name(), err)

			return nil
		}

		db.addVulnerability(pa, pkgNames)

		return nil
	})

	if errored {
		_, _ = fmt.Fprintf(os.Stderr, "\n")
	}

	if err != nil {
		return fmt.Errorf("could not read OSV database directory: %w", err)
	}

	return nil
}

func NewDirDB(config Config, offline bool, pkgNames []string) (*DirDB, error) {
	db := &DirDB{
		name:             config.Name,
		identifier:       config.Identifier(),
		LocalPath:        config.URL,
		WorkingDirectory: config.WorkingDirectory,
		Offline:          offline,
	}
	if err := db.load(pkgNames); err != nil {
		return nil, fmt.Errorf("unable to load OSV database: %w", err)
	}

	return db, nil
}
