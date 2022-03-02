package semver

import (
	"fmt"
	"strings"
)

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
