package database

import (
	"encoding/json"
	"fmt"
)

type Vulnerabilities []OSV

func (vs Vulnerabilities) Unique() Vulnerabilities {
	var vulnerabilities Vulnerabilities

	for _, vulnerability := range vs {
		if !vulnerabilities.Includes(vulnerability) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}

func (vs Vulnerabilities) Includes(vulnerability OSV) bool {
	for _, osv := range vs {
		if osv.ID == vulnerability.ID {
			return true
		}

		if osv.isAliasOf(vulnerability) {
			return true
		}
		if vulnerability.isAliasOf(osv) {
			return true
		}
	}

	return false
}

// MarshalJSON ensures that if there are no vulnerabilities,
// an empty array is used as the value instead of "null"
func (vs Vulnerabilities) MarshalJSON() ([]byte, error) {
	if len(vs) == 0 {
		return []byte("[]"), nil
	}

	type innerVulnerabilities Vulnerabilities

	out, err := json.Marshal(innerVulnerabilities(vs))

	if err != nil {
		return out, fmt.Errorf("%w", err)
	}

	return out, nil
}
