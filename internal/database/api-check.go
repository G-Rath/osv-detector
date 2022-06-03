package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"osv-detector/internal"
	"path"
)

func (db APIDB) buildAPIPayload(pkg internal.PackageDetails) apiPayload {
	var payload apiPayload

	if pkg.Commit == "" {
		payload.Package.Name = pkg.Name
		payload.Package.Ecosystem = pkg.Ecosystem
		payload.Version = pkg.Version
	} else {
		payload.Commit = pkg.Commit
	}

	return payload
}

func (db APIDB) bulkEndpoint() string {
	u := *db.BaseURL

	u.Path = path.Join(u.Path, "querybatch")

	return u.String()
}

type ObjectWithID struct {
	ID string `json:"id"`
}

var ErrAPICouldNotMarshalPayload = errors.New("could not marshal payload")
var ErrAPIRequestInvalid = errors.New("api request invalid")
var ErrAPIRequestFailed = errors.New("api request failed")
var ErrAPIUnexpectedResponse = errors.New("api returned unexpected status")
var ErrAPIUnreadableResponse = errors.New("could not read response body")
var ErrAPIResponseNotJSON = errors.New("api response could not be parsed as json")
var ErrAPIResultsCountMismatch = errors.New("api results count mismatch")

func (db APIDB) checkBatch(pkgs []internal.PackageDetails) ([]VulnsOrError, error) {
	payloads := make([]apiPayload, 0, len(pkgs))

	for _, pkg := range pkgs {
		payloads = append(payloads, db.buildAPIPayload(pkg))
	}

	jsonData, err := json.Marshal(struct {
		Queries []apiPayload `json:"queries"`
	}{payloads})

	if err != nil {
		return []VulnsOrError{}, fmt.Errorf("%v: %w", ErrAPICouldNotMarshalPayload, err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		db.bulkEndpoint(),
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return []VulnsOrError{}, fmt.Errorf("%v: %w", ErrAPIRequestInvalid, err)
	}

	// fmt.Println(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []VulnsOrError{}, fmt.Errorf("%v: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	var body []byte

	body, err = io.ReadAll(resp.Body)

	// fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return []VulnsOrError{}, fmt.Errorf("%w (%d)", ErrAPIUnexpectedResponse, resp.StatusCode)
	}

	// body, err = io.ReadAll(resp.Body)

	// fmt.Println(string(body))

	if err != nil {
		return []VulnsOrError{}, fmt.Errorf("%v: %w", ErrAPIUnreadableResponse, err)
	}

	var parsed struct {
		Results []struct {
			Vulns []ObjectWithID `json:"vulns"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &parsed)

	if err != nil {
		return []VulnsOrError{}, fmt.Errorf("%v: %w", ErrAPIResponseNotJSON, err)
	}

	vulnerabilities := make([]VulnsOrError, 0, len(parsed.Results))

	for i, r := range parsed.Results {
		vulnerabilities = append(vulnerabilities, VulnsOrError{
			Index: i,
			Vulns: db.expand(r.Vulns),
			Err:   nil,
		})
	}

	if len(pkgs) != len(vulnerabilities) {
		return vulnerabilities, fmt.Errorf(
			"%w - expected to get %d but got %d",
			ErrAPIResultsCountMismatch,
			len(pkgs),
			len(vulnerabilities),
		)
	}

	return vulnerabilities, nil
}

func (db APIDB) expand(ids []ObjectWithID) []OSV {
	osvs := make([]OSV, 0, len(ids))

	for _, withID := range ids {
		// if we error, still report the vulnerability as hopefully the ID should be
		// enough to manually look up the details - in future we should ideally warn
		// the user too, but for now we just silently eat the error
		osv, _ := db.Fetch(withID.ID)
		osvs = append(osvs, osv)
	}

	return osvs
}

func batchPkgs(pkgs []internal.PackageDetails, batchSize int) [][]internal.PackageDetails {
	batches := make(
		[][]internal.PackageDetails,
		0,
		(len(pkgs)/batchSize)+int(math.Min(float64(len(pkgs)%batchSize), 1)),
	)

	for i := 0; i < len(pkgs); i += batchSize {
		end := i + batchSize

		if end > len(pkgs) {
			end = len(pkgs)
		}

		batches = append(batches, pkgs[i:end])
	}

	return batches
}

func (db APIDB) Check(pkgs []internal.PackageDetails) ([]VulnsOrError, error) {
	batches := batchPkgs(pkgs, db.BatchSize)

	vulnerabilities := make([]VulnsOrError, 0, len(pkgs))

	for _, batch := range batches {
		results, err := db.checkBatch(batch)

		if err != nil {
			return nil, err
		}

		vulnerabilities = append(vulnerabilities, results...)
	}

	return vulnerabilities, nil
}
