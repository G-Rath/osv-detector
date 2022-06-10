package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

func (db APIDB) osvEndpoint(id string) string {
	u := *db.BaseURL

	u.Path = path.Join(u.Path, "vulns", id)

	return u.String()
}

// Fetch gets the details of a specific OSV from the osv.dev database
func (db APIDB) Fetch(id string) (OSV, error) {
	var osv OSV

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		db.osvEndpoint(id),
		http.NoBody,
	)

	if err != nil {
		return osv, fmt.Errorf("%v: %w", ErrAPIRequestInvalid, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return osv, fmt.Errorf("%v: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	var body []byte

	body, err = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return osv, fmt.Errorf(
			"%w (%s %s %d)",
			ErrAPIUnexpectedResponse,
			resp.Request.Method,
			resp.Request.URL,
			resp.StatusCode,
		)
	}

	if err != nil {
		return osv, fmt.Errorf(
			"%v (%s %s): %w",
			ErrAPIUnreadableResponse,
			resp.Request.Method,
			resp.Request.URL,
			err,
		)
	}

	err = json.Unmarshal(body, &osv)

	if err != nil {
		return osv, fmt.Errorf(
			"%v (%s %s): %w",
			ErrAPIResponseNotJSON,
			resp.Request.Method,
			resp.Request.URL,
			err,
		)
	}

	return osv, nil
}
