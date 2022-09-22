package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

type Components []*big.Int

type Version struct {
	LeadingV   bool
	Components Components
	Build      string
}

func (components *Components) Fetch(n int) *big.Int {
	if len(*components) <= n {
		return big.NewInt(0)
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
