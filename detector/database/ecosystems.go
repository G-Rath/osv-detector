package database

import (
	"fmt"
	"sort"
)

func toSliceOfEcosystems(ecosystemsMap map[Ecosystem]struct{}) []Ecosystem {
	ecosystems := make([]Ecosystem, 0, len(ecosystemsMap))

	for ecosystem := range ecosystemsMap {
		ecosystems = append(ecosystems, ecosystem)
	}

	return ecosystems
}

func (db *OSVDatabase) ListEcosystems() []Ecosystem {
	ecosystems := make(map[Ecosystem]struct{})

	for _, vulnerability := range db.Vulnerabilities(true) {
		if vulnerability.Affected == nil {
			fmt.Printf("Skipping %s as it does not have an 'affected' property", vulnerability.ID)

			continue
		}

		for _, affected := range vulnerability.Affected {
			ecosystems[affected.Package.Ecosystem] = struct{}{}
		}
	}

	slicedEcosystems := toSliceOfEcosystems(ecosystems)

	sort.Slice(slicedEcosystems, func(i, j int) bool {
		return slicedEcosystems[i] < slicedEcosystems[j]
	})

	return slicedEcosystems
}

func (db *OSVDatabase) ListEcosystemVulnerabilities(ecosystem Ecosystem) []OSV {
	var vulnerabilities []OSV

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.Affected == nil {
			fmt.Printf("Skipping %s as it does not have an 'affected' property", vulnerability.ID)

			continue
		}

		if vulnerability.AffectsEcosystem(ecosystem) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}
