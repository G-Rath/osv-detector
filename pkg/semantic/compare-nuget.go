package semantic

import "strings"

type NuGetVersion struct {
	SemverLikeVersion
}

func (v NuGetVersion) Compare(w NuGetVersion) int {
	componentDiff := compareNumericComponents(v.Components, w.Components)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuildComponents(strings.ToLower(v.Build), strings.ToLower(w.Build))
}

func (v NuGetVersion) CompareStr(str string) int {
	return v.Compare(parseNuGetVersion(str))
}

func parseNuGetVersion(str string) NuGetVersion {
	return NuGetVersion{ParseSemverLikeVersion(str, 4)}
}
