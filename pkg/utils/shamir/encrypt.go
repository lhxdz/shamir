package shamir

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/log"
	math2 "shamir/pkg/utils/math"
)

const (
	MinThreshold = 2
)

var minPrime = math2.NextPrime(big.NewInt(math.MaxInt64))

// Encrypt shamir加密算法，输入秘密secret，门限值threshold，密钥个数keysNumber
// 返回加密结果密钥对keys，和加密使用的素数，都需谨慎保存
func Encrypt(secret *big.Int, threshold, keysNumber int, fast bool) (keys []code.Key, prime *big.Int, err error) {
	if err = encryptCheck(secret, threshold, keysNumber); err != nil {
		return nil, nil, err
	}

	return encrypt(secret, threshold, keysNumber, fast)
}

// CompoundEncrypt 用于加密复合型秘密，由于秘密太大，拆分成多个子秘密加密，将会生成复合型密钥和复合型素数
func CompoundEncrypt(secret []*big.Int, threshold, keysNumber int, fast bool) (keys []code.CompoundKey, prime []*big.Int, err error) {
	if err = CompoundEncryptCheck(secret, threshold, keysNumber); err != nil {
		return nil, nil, err
	}

	return compoundEncrypt(secret, threshold, keysNumber, fast)
}

// HashEncrypt 将秘密进行hash计算，并将计算结果一并加密进入密钥中
func HashEncrypt(secret []*big.Int, threshold, keysNumber int, fast bool) (keys []code.CompoundKey, prime []*big.Int, err error) {
	if err = CompoundEncryptCheck(secret, threshold, keysNumber); err != nil {
		return nil, nil, err
	}

	hash := sha256.New()
	for i := 0; i < len(secret); i++ {
		_, err = hash.Write(secret[i].Bytes())
		if err != nil {
			log.Errorf("hash check sum failed: %v", err)
			return nil, nil, HashCheckFailed
		}
	}

	hashInt := code.EncodeSecret(string(hash.Sum(nil)))
	// hash值的加密不使用fast，实时计算prime
	hashKeys, hashPrime, err := encrypt(hashInt, threshold, keysNumber, false)
	if err != nil {
		log.Errorf("encrypt hash value failed: %v", err)
		return nil, nil, err
	}

	keys, prime, err = compoundEncrypt(secret, threshold, keysNumber, fast)
	if err != nil {
		log.Errorf("hash encrypt failed: %v", err)
		return nil, nil, err
	}

	prime = append(prime, hashPrime)
	for i, tmpKey := range hashKeys {
		keys[i].X = append(keys[i].X, tmpKey.X)
		keys[i].Y = append(keys[i].Y, tmpKey.Y)
	}

	return
}

// private

func encrypt(secret *big.Int, threshold, keysNumber int, fast bool) (keys []code.Key, prime *big.Int, err error) {
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

	xKeys, err := math2.NewRandGenerator(minInt(minPrime, prime)).RandIntListNoRepeat(keysNumber)
	if err != nil {
		return nil, nil, err
	}

	keys = make([]code.Key, 0, keysNumber)
	for _, xKey := range xKeys {
		yKey := compute(coefficients, prime, xKey)
		keys = append(keys, code.Key{X: xKey, Y: yKey})
	}

	return
}

func compoundEncrypt(secret []*big.Int, threshold, keysNumber int, fast bool) (keys []code.CompoundKey, prime []*big.Int, err error) {
	keys = make([]code.CompoundKey, keysNumber, keysNumber)
	prime = make([]*big.Int, 0, len(secret))

	for _, tmpSecret := range secret {
		tmpKeys, tmpPrime, e := encrypt(tmpSecret, threshold, keysNumber, fast)
		if e != nil {
			log.Errorf("compound encrypt failed: %v", e)
			return nil, nil, e
		}

		prime = append(prime, tmpPrime)
		for i, tmpKey := range tmpKeys {
			keys[i].X = append(keys[i].X, tmpKey.X)
			keys[i].Y = append(keys[i].Y, tmpKey.Y)
		}
	}

	return
}

func encryptCheck(secret *big.Int, threshold, keysNumber int) error {
	if secret == nil {
		return fmt.Errorf("nil point of secret")
	}

	return tnCheck(threshold, keysNumber)
}

func CompoundEncryptCheck(secret []*big.Int, threshold, keysNumber int) error {
	for _, tmpSecret := range secret {
		if tmpSecret == nil {
			return fmt.Errorf("nil point of secret")
		}
	}

	return tnCheck(threshold, keysNumber)
}

func tnCheck(threshold, keysNumber int) error {
	if threshold > keysNumber {
		return fmt.Errorf("threshold(%d) can not bigger than keys number(%d)", threshold, keysNumber)
	}
	if threshold < MinThreshold {
		return fmt.Errorf("threshold(%d) can not smaller than %d", threshold, MinThreshold)
	}

	return nil
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
