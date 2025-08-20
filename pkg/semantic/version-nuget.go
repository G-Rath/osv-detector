package semantic

import "strings"

type NuGetVersion struct {
	SemverLikeVersion
}

func (v NuGetVersion) Compare(w NuGetVersion) int {
	if diff := v.Components.Cmp(w.Components); diff != 0 {
		return diff
	}

	return compareBuildComponents(strings.ToLower(v.Build), strings.ToLower(w.Build))
}

func (v NuGetVersion) CompareStr(str string) (int, error) {
	return v.Compare(parseNuGetVersion(str)), nil
}

func parseNuGetVersion(str string) NuGetVersion {
	return NuGetVersion{ParseSemverLikeVersion(str, 4)}
}
