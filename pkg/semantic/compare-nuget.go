package semantic

import "strings"

func compareForNuGet(v, w Version) int {
	vComponents, vBuild := v.fetchComponentsAndBuild(4)
	wComponents, wBuild := w.fetchComponentsAndBuild(4)

	componentDiff := compareNumericComponents(vComponents, wComponents)

	if componentDiff != 0 {
		return componentDiff
	}

	return compareBuildComponents(strings.ToLower(vBuild), strings.ToLower(wBuild))
}
