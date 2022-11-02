package shamir

import (
	"fmt"
	"math"
	"math/big"
	math2 "shamir/pkg/utils/math"
)

const (
	minThreshold = 2
)

var minPrime = math2.NextPrime(big.NewInt(math.MaxInt64))

// Encrypt shamir加密算法，输入秘密secret，门限值threshold，密钥个数keysNumber
// 返回加密结果密钥对keys，和加密使用的素数，都需谨慎保存
func Encrypt(secret *big.Int, threshold, keysNumber int, fast bool) (keys []Key, prime *big.Int, err error) {
	if threshold > keysNumber {
		return nil, nil, fmt.Errorf("threshold(%d) can not bigger than keys number(%d)", threshold, keysNumber)
	}
	if threshold < minThreshold {
		return nil, nil, fmt.Errorf("threshold(%d) can not smaller than %d", threshold, minThreshold)
	}

	if minPrime.Cmp(secret) > 0 {
		// 若秘密太小，则使用默认质数
		prime = new(big.Int).Set(minPrime)
	} else {
		if fast {
			prime = math2.FastPrime(secret)
		} else {
			prime = math2.NextPrime(secret)
		}
	}

	coefficients := make([]*big.Int, 0, threshold)
	// secret作为系数a0
	coefficients = append(coefficients, secret)
	tmpCoefficients, err := math2.NewRandGenerator(prime).RandIntList(threshold - 1)
	if err != nil {
		return nil, nil, err
	}
	coefficients = append(coefficients, tmpCoefficients...)

	xKeys, err := math2.NewRandGenerator(minInt(minPrime, prime)).RandIntList(keysNumber)
	if err != nil {
		return nil, nil, err
	}

	keys = make([]Key, 0, keysNumber)
	for _, xKey := range xKeys {
		yKey := compute(coefficients, prime, xKey)
		keys = append(keys, Key{X: xKey, Y: yKey})
	}

	return
}

// 计算 f(x) = (a0 + a1*(x^1) + a2*(x^2) + ... an*(x^n)) mod prime
func compute(coefficients []*big.Int, prime, x *big.Int) *big.Int {
	y := big.NewInt(0)
	for i, coefficient := range coefficients {
		// x**i mod prime
		tmp := new(big.Int).Exp(x, big.NewInt(int64(i)), prime)
		// ai*tmp
		tmp = tmp.Mul(tmp, coefficient)
		y = new(big.Int).Mod(y.Add(y, tmp), prime)
	}
	return y
}

func minInt(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}
