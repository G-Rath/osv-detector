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

func (db APIDB) checkBatch(pkgs []internal.PackageDetails) ([][]ObjectWithID, error) {
	payloads := make([]apiPayload, 0, len(pkgs))

	for _, pkg := range pkgs {
		payloads = append(payloads, db.buildAPIPayload(pkg))
	}

	jsonData, err := json.Marshal(struct {
		Queries []apiPayload `json:"queries"`
	}{payloads})

	if err != nil {
		return [][]ObjectWithID{}, fmt.Errorf("%v: %w", ErrAPICouldNotMarshalPayload, err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		db.bulkEndpoint(),
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return [][]ObjectWithID{}, fmt.Errorf("%v: %w", ErrAPIRequestInvalid, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return [][]ObjectWithID{}, fmt.Errorf("%v: %w", ErrAPIRequestFailed, err)
	}

	defer resp.Body.Close()

	var body []byte

	body, err = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return [][]ObjectWithID{}, fmt.Errorf("%w (%d)", ErrAPIUnexpectedResponse, resp.StatusCode)
	}

	if err != nil {
		return [][]ObjectWithID{}, fmt.Errorf("%v: %w", ErrAPIUnreadableResponse, err)
	}

	var parsed struct {
		Results []struct {
			Vulns []ObjectWithID `json:"vulns"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &parsed)

	if err != nil {
		return [][]ObjectWithID{}, fmt.Errorf("%v: %w", ErrAPIResponseNotJSON, err)
	}

	vulnerabilities := make([][]ObjectWithID, 0, len(parsed.Results))

	for _, r := range parsed.Results {
		vulnerabilities = append(vulnerabilities, r.Vulns)
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

func findOrDefault(vulns Vulnerabilities, def OSV) OSV {
	for _, vuln := range vulns {
		if vuln.ID == def.ID {
			return vuln
		}
	}

	return def
}

func (db APIDB) Check(pkgs []internal.PackageDetails) ([]Vulnerabilities, error) {
	batches := batchPkgs(pkgs, db.BatchSize)

	vulnerabilities := make([]Vulnerabilities, 0, len(pkgs))

	for _, batch := range batches {
		results, err := db.checkBatch(batch)

		if err != nil {
			return nil, err
		}

		for _, withIDS := range results {
			vulns := make(Vulnerabilities, 0, len(withIDS))

			for _, withID := range withIDS {
				vulns = append(vulns, OSV{ID: withID.ID})
			}

			vulnerabilities = append(vulnerabilities, vulns)
		}
	}

	var osvs Vulnerabilities

	for _, vulns := range vulnerabilities {
		osvs = append(osvs, vulns...)
	}

	osvs = osvs.Unique()

	ids := make([]string, 0, len(osvs))

	for _, osv := range osvs {
		ids = append(ids, osv.ID)
	}

	osvs = db.FetchAll(ids)

	for _, vulns := range vulnerabilities {
		for i := range vulns {
			vulns[i] = findOrDefault(osvs, vulns[i])
		}
	}

	return vulnerabilities, nil
}
