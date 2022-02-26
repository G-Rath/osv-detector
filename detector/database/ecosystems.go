package database

import "fmt"

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

	return toSliceOfEcosystems(ecosystems)
}
