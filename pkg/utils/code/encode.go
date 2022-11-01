package code

import (
	"math/big"
	"strings"
)

// EncodeSecret 将字符串秘密编码成为大整数，方便加密
func EncodeSecret(secret string) *big.Int {
	return new(big.Int).SetBytes([]byte(secret))
}

// EncodeKey 将密钥字符串恢复成大整数密钥
func EncodeKey(key string) (*big.Int, bool) {
	return new(big.Int).SetString(key, base)
}

// EncodeKeys 将密钥字符串恢复成大整数密钥链
func EncodeKeys(key string) ([]*big.Int, bool) {
	subKeys := strings.Split(key, splitKey)
	bigIntegers := make([]*big.Int, 0, len(subKeys))
	for _, subKey := range subKeys {
		subInteger, ok := EncodeKey(subKey)
		if !ok {
			return nil, false
		}
		bigIntegers = append(bigIntegers, subInteger)
	}
	return bigIntegers, true
}
