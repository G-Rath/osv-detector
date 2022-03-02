package semver

import (
	"fmt"
	"strings"
)

// type Version struct {
// 	Major int
// 	Minor int
// 	Patch int
// 	Other int
// 	Build string
// }
//
// func (v Version) ToString() string {
// 	return fmt.Sprintf(
// 		"%d.%d.%d.%d-%s",
// 		v.Major,
// 		v.Minor,
// 		v.Patch,
// 		v.Other,
// 		v.Build,
// 	)
// }

type Version struct {
	Components []int
	Build      string
}

func (v Version) ToString() string {
	str := ""

	for _, component := range v.Components {
		str += fmt.Sprintf("%d.", component)
	}

	str = strings.TrimSuffix(str, ".")
	str += v.Build

	return str
}
