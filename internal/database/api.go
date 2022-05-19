package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"osv-detector/internal"
)

type APIOSVDatabase struct {
}

type apiPayload struct {
	Commit  string `json:"commit,omitempty"`
	Version string `json:"version,omitempty"`
	Package struct {
		Name      string             `json:"name"`
		Ecosystem internal.Ecosystem `json:"ecosystem"`
	} `json:"package,omitempty"`
}

func (db APIOSVDatabase) buildAPIPayload(pkg internal.PackageDetails) apiPayload {
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

func (db APIOSVDatabase) Check(pkgs []internal.PackageDetails) []Vulnerabilities {
	return db.VulnerabilitiesAffectingPackages(pkgs)
}

func (db APIOSVDatabase) VulnerabilitiesAffectingPackages(pkgs []internal.PackageDetails) []Vulnerabilities {
	payloads := make([]apiPayload, 0, len(pkgs))

	for _, pkg := range pkgs {
		payloads = append(payloads, db.buildAPIPayload(pkg))
	}

	jsonData, err := json.Marshal(struct {
		Queries []apiPayload `json:"queries"`
	}{payloads})

	if err != nil {
		fmt.Printf("error marshaling payload: %v", err)

		return []Vulnerabilities{}
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api-staging.osv.dev/v1/querybatch",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Printf("error building request: %v", err)

		return []Vulnerabilities{}
	}

	// fmt.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error making request: %v", err)

		return []Vulnerabilities{}
	}

	defer resp.Body.Close()

	var body []byte

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response was not 200: %d", resp.StatusCode)

		return []Vulnerabilities{}
	}

	body, err = io.ReadAll(resp.Body)

	// fmt.Println(string(body))

	if err != nil {
		fmt.Printf("error reading response body: %v", err)

		return []Vulnerabilities{}
	}

	var parsed struct {
		Results []struct {
			Vulns Vulnerabilities `json:"vulns"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &parsed)

	if err != nil {
		fmt.Printf("error reading response body: %v", err)

		return []Vulnerabilities{}
	}

	vulnerabilities := make([]Vulnerabilities, 0, len(parsed.Results))

	for _, r := range parsed.Results {
		vulnerabilities = append(vulnerabilities, r.Vulns)
	}

	return vulnerabilities
}

func (db APIOSVDatabase) VulnerabilitiesAffectingPackage(pkg internal.PackageDetails) Vulnerabilities {
	var vulnerabilities Vulnerabilities

	jsonData, err := json.Marshal(db.buildAPIPayload(pkg))

	if err != nil {
		fmt.Printf("error marshaling payload: %v", err)

		return vulnerabilities
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://api.osv.dev/v1/query",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Printf("error building request: %v", err)

		return vulnerabilities
	}

	// fmt.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error making request: %v", err)

		return vulnerabilities
	}

	defer resp.Body.Close()

	var body []byte

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response was not 200: %d", resp.StatusCode)

		return vulnerabilities
	}

	body, err = io.ReadAll(resp.Body)

	// fmt.Println(string(body))

	if err != nil {
		fmt.Printf("error reading response body: %v", err)

		return vulnerabilities
	}

	var parsed struct {
		Vulns Vulnerabilities `json:"vulns"`
	}

	err = json.Unmarshal(body, &parsed)

	if err != nil {
		fmt.Printf("error reading response body: %v", err)

		return vulnerabilities
	}

	return parsed.Vulns
}
