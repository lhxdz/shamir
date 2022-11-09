package shamir

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/log"
)

const (
	negative = -1
)

var HashCheckFailed = errors.New("secret hash check failed")

// Decrypt 根据密钥对解密，传入 len(keys) 必须是门限值threshold, prime 也必须是加密时使用的素数。
// keys要求不能有重复元素，否则panic
func Decrypt(keys []code.Key, prime *big.Int) (secret *big.Int, err error) {
	if err = decryptCheck(keys, prime); err != nil {
		return nil, err
	}

	return decrypt(keys, prime), nil
}

// CompoundDecrypt 用于解密复合型秘密，使用复合型密钥和复合型素数，解密出复合型秘密
func CompoundDecrypt(keys []code.CompoundKey, prime []*big.Int) (secret []*big.Int, err error) {
	if err = checkCompoundKeys(keys, prime); err != nil {
		return nil, err
	}

	return compoundDecrypt(keys, prime), nil
}

// HashDecrypt 有hash校验的解密，将会从密钥中解密出秘密和hash值，
// 若未通过hash校验返回错误
// 输入的复合key中复合的最后一对key作为hash校验的依据
func HashDecrypt(keys []code.CompoundKey, prime []*big.Int) (secret []*big.Int, err error) {
	if err := checkCompoundKeys(keys, prime); err != nil {
		return nil, err
	}

	// 因为最后一个复合key是用作校验的hash值，所以长度必须大于1
	if len(prime) < 2 {
		return nil, fmt.Errorf("invalid input keys")
	}

	secret = compoundDecrypt(keys, prime)
	checkHash := code.DecodeSecret(secret[len(secret)-1])
	secret = secret[:len(secret)-1]

	newHash := sha256.New()
	for i := 0; i < len(secret); i++ {
		_, err := newHash.Write(secret[i].Bytes())
		if err != nil {
			log.Errorf("hash check sum failed: %v", err)
			return nil, HashCheckFailed
		}
	}

	if string(newHash.Sum(nil)) != checkHash {
		return nil, HashCheckFailed
	}

	return
}

// private

func decrypt(keys []code.Key, prime *big.Int) *big.Int {
	result := big.NewInt(0)
	for i, key := range keys {
		tmp := product(getXKeysExceptI(keys, i), key.X, prime)
		tmp = tmp.Mul(tmp, key.Y)
		result = result.Add(result, tmp)
		result = result.Mod(result, prime)
	}
	return result
}

func compoundDecrypt(keys []code.CompoundKey, prime []*big.Int) []*big.Int {
	result := make([]*big.Int, 0, len(prime))
	for i, tmpPrime := range prime {
		tmpKeys := make([]code.Key, 0, len(keys))
		for _, key := range keys {
			tmpKeys = append(tmpKeys, code.Key{X: key.X[i], Y: key.Y[i]})
		}

		result = append(result, decrypt(tmpKeys, tmpPrime))
	}
	return result
}

func decryptCheck(keys []code.Key, prime *big.Int) error {
	if prime == nil {
		return fmt.Errorf("invalid nil point prime")
	}

	for _, key := range keys {
		if key.X == nil || key.Y == nil {
			return fmt.Errorf("invalid nil point key")
		}
	}

	return nil
}

func checkCompoundKeys(keys []code.CompoundKey, prime []*big.Int) error {
	if len(keys) == 0 {
		return nil
	}
	if len(keys[0].X) != len(prime) {
		return fmt.Errorf("input key's compound count not equal prime's count")
	}

	for _, key := range keys {
		if len(key.X) != len(key.Y) {
			return fmt.Errorf("input key_x's compound count not equal key_y's count")
		}

		for i := range key.X {
			if key.X[i] == nil || key.Y[i] == nil {
				return fmt.Errorf("invalid nil point key")
			}
		}
	}

	return nil
}

func getXKeysExceptI(keys []code.Key, i int) []*big.Int {
	xKeys := make([]*big.Int, 0, len(keys)-1)
	for j, key := range keys {
		if j == i {
			continue
		}
		xKeys = append(xKeys, key.X)
	}

	return xKeys
}

// 求((-xKeys[0])*...(-xKeys[n])) / ((xKeysI - xKeys[0])*...(xKeysI - xKeys[n])) mod prime
// 其中 n = (len(xKeys)-1)
// 其中的xKey不能为0
// prime为0会panic
// 若其中xKeys某一项和xKeysI值相等，则会panic
func product(xKeys []*big.Int, xKeysI *big.Int, prime *big.Int) *big.Int {
	result := big.NewInt(1)
	if len(xKeys)%2 != 0 {
		// 奇数个负数分子相乘为奇数
		result = big.NewInt(-1)
	}
	for _, key := range xKeys {
		result = result.Mul(result, key)
		result = result.Mod(result, prime)
		denominator := new(big.Int).Sub(xKeysI, key)
		if denominator.Sign() == negative {
			// 若分母为负数,则应先求模，再求 x**y mod p
			denominator = denominator.Mod(denominator, prime)
		}
		denominator = denominator.Exp(denominator, big.NewInt(-1), prime)
		// 分子分母取模结果相乘再取模
		result = result.Mul(result, denominator)
		result = result.Mod(result, prime)
	}

	return result
}
