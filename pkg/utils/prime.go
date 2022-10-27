package utils

import (
	"math/big"
)

const (
	defaultTestTime = 10
)

var incrementalSize = big.NewInt(1)

// NextPrime 返回 >= num 的第一个素数
func NextPrime(num *big.Int) *big.Int {
	if num.Cmp(big.NewInt(3)) < 0 {
		return big.NewInt(2)
	}

	result := new(big.Int).Set(num)

	for !result.ProbablyPrime(defaultTestTime) {
		result.Add(result, incrementalSize)
	}

	return result
}
