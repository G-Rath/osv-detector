package semantic

import (
	"fmt"
	"strings"
)

type Components []int

type Version struct {
	LeadingV   bool
	Components Components
	Build      string
}

func (components *Components) Fetch(n int) int {
	if len(*components) <= n {
		return 0
	}

	return (*components)[n]
}

func (v *Version) String() string {
	str := ""

	if v.LeadingV {
		str += "v"
	}

	for _, component := range v.Components {
		str += fmt.Sprintf("%d.", component)
	}

	str = strings.TrimSuffix(str, ".")
	str += v.Build

	return str
}
