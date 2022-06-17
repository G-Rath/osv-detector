package database

import (
	"errors"
	"fmt"
	"net/url"
	"osv-detector/internal"
)

type APIDB struct {
	name       string
	identifier string
	BaseURL    *url.URL
	BatchSize  int
}

func (db APIDB) Identifier() string {
	return db.identifier
}

func (db APIDB) Name() string {
	return db.name
}

type apiQuery struct {
	Commit  string `json:"commit,omitempty"`
	Version string `json:"version,omitempty"`
	Package struct {
		Name      string             `json:"name"`
		Ecosystem internal.Ecosystem `json:"ecosystem"`
	} `json:"package,omitempty"`
}

var ErrOfflineDatabaseNotSupported = errors.New("API database does not support being used offline")
var ErrInvalidBatchSize = errors.New("batch size must be greater than 0")

func NewAPIDB(config Config, offline bool, batchSize int) (*APIDB, error) {
	if offline {
		return nil, ErrOfflineDatabaseNotSupported
	}

	if batchSize < 1 {
		return nil, ErrInvalidBatchSize
	}

	u, err := url.ParseRequestURI(config.URL)

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid url: %w", config.URL, err)
	}

	return &APIDB{
		name:       config.Name,
		identifier: config.Identifier(),
		BaseURL:    u,
		BatchSize:  batchSize,
	}, nil
}
