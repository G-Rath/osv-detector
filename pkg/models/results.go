package models

// VulnerabilityResults represents the vulnerabilities results for each package source
type VulnerabilityResults struct {
	Results []PackageSource `json:"results"`
}
