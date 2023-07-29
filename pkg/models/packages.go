package models

import "sort"

type PackageSource struct {
	FilePath string   `json:"filePath"`
	ParsedAs string   `json:"parsedAs"`
	Packages Packages `json:"packages"`
}

type PackageInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	Ecosystem Ecosystem `json:"ecosystem"`
	CompareAs Ecosystem `json:"compareAs"`
}

type PackageInfoWithVulnerabilities struct {
	PackageInfo

	Vulnerabilities Vulnerabilities `json:"vulnerabilities"`
	Ignored         Vulnerabilities `json:"ignored"`
}

type Packages []PackageInfo

func toSliceOfEcosystems(ecosystemsMap map[Ecosystem]struct{}) []Ecosystem {
	ecosystems := make([]Ecosystem, 0, len(ecosystemsMap))

	for ecosystem := range ecosystemsMap {
		if ecosystem == "" {
			continue
		}

		ecosystems = append(ecosystems, ecosystem)
	}

	return ecosystems
}

func (ps Packages) Ecosystems() []Ecosystem {
	ecosystems := make(map[Ecosystem]struct{})

	for _, pkg := range ps {
		ecosystems[pkg.Ecosystem] = struct{}{}
	}

	slicedEcosystems := toSliceOfEcosystems(ecosystems)

	sort.Slice(slicedEcosystems, func(i, j int) bool {
		return slicedEcosystems[i] < slicedEcosystems[j]
	})

	return slicedEcosystems
}
