package math

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type randGenerator struct {
	limit *big.Int
}

// NewRandGenerator 获得一个随机数生成器，将生成(0, max)区间的随机数，不包含0和max
func NewRandGenerator(max *big.Int) *randGenerator {
	newMax := new(big.Int).Set(max)
	newMax = newMax.Sub(newMax, incrementalSize)
	return &randGenerator{
		limit: newMax,
	}
}

func (r *randGenerator) RandInt() (*big.Int, error) {
	random, err := rand.Int(rand.Reader, r.limit)
	if err != nil {
		return nil, fmt.Errorf("get random bit int failed: %w", err)
	}
	return random.Add(random, incrementalSize), nil
}

// RandIntList 返回一个随机数列表，其中数可以不重复，且范围在(0, max)区间
func (r *randGenerator) RandIntList(num int) ([]*big.Int, error) {
	return r.randIntList(num, false)
}

// RandIntListNoRepeat 返回一个随机数列表，其中数不重复，且范围在(0, max)区间
func (r *randGenerator) RandIntListNoRepeat(num int) ([]*big.Int, error) {
	return r.randIntList(num, true)
}

func (r *randGenerator) randIntList(num int, notRepeat bool) ([]*big.Int, error) {
	if notRepeat && big.NewInt(int64(num)).Cmp(r.limit) > 0 {
		return nil, fmt.Errorf("number out of range(0, %d)", r.limit.Int64()+1)
	}

	result := make([]*big.Int, 0, num)
	for len(result) < num {
		random, err := r.RandInt()
		if err != nil {
			return nil, err
		}

		if notRepeat && InList(result, random) {
			continue
		}
		result = append(result, random)
	}

	return result, nil
}

func InList(list []*big.Int, num *big.Int) bool {
	for _, key := range list {
		if num.Cmp(key) == 0 {
			return true
		}
	}

	return false
}
