package database

import (
	"errors"
	"fmt"
	"net/url"
	"osv-detector/internal"
)

type APIDB struct {
	BaseURL *url.URL
}

type apiPayload struct {
	Commit  string `json:"commit,omitempty"`
	Version string `json:"version,omitempty"`
	Package struct {
		Name      string             `json:"name"`
		Ecosystem internal.Ecosystem `json:"ecosystem"`
	} `json:"package,omitempty"`
}

var ErrOfflineDatabaseNotSupported = errors.New("API database does not support being used offline")

func NewAPIDB(baseURL string, offline bool) (*APIDB, error) {
	if offline {
		return nil, ErrOfflineDatabaseNotSupported
	}

	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid url: %w", baseURL, err)
	}

	return &APIDB{BaseURL: u}, nil
}
