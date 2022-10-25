package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type randGenerator struct {
	limit *big.Int
}

// 获得一个随机数生成器，将生成(0, max)区间的随机数，不包含0和max
func newRandGenerator(max *big.Int) *randGenerator {
	newMax := new(big.Int).Set(max)
	newMax = newMax.Sub(newMax, big.NewInt(1))
	return &randGenerator{
		limit: newMax,
	}
}

func (r *randGenerator) randInt() (*big.Int, error) {
	random, err := rand.Int(rand.Reader, r.limit)
	if err != nil {
		return nil, fmt.Errorf("get random bit int failed: %w", err)
	}
	return random.Add(random, incrementalSize), nil
}

func (r *randGenerator) randIntList(num int) ([]*big.Int, error) {
	result := make([]*big.Int, num)
	for i := 0; i < num; i++ {
		random, err := r.randInt()
		if err != nil {
			return nil, err
		}
		result = append(result, random)
	}

	return result, nil
}
