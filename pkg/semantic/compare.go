package semantic

import (
	"math/big"
)

func convertToBigInt(str string) (*big.Int, bool) {
	i, ok := new(big.Int).SetString(str, 10)

	return i, ok
}

func minInt(x, y int) int {
	if x > y {
		return y
	}

	return x
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}

	return x
}

func compareComponents(a Components, b Components) int {
	numberOfComponents := maxInt(len(a), len(b))

	for i := 0; i < numberOfComponents; i++ {
		diff := a.Fetch(i).Cmp(b.Fetch(i))

		if diff != 0 {
			return diff
		}
	}

	return 0
}
