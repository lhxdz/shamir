package code

import (
	"math/big"
	"strings"

	"shamir/pkg/utils/compute"
)

// EncodeSecret 将字符串秘密编码成为大整数，方便加密
func EncodeSecret(secret string) *big.Int {
	return new(big.Int).SetBytes([]byte(secret))
}

func EncodeCompoundSecret(secret string, splitLen int) []*big.Int {
	result := make([]*big.Int, 0, getBucketCounts(len(secret), splitLen))
	for i := 0; i < len(secret); i += splitLen {
		result = append(result, EncodeSecret(secret[i:compute.Min(len(secret), i+splitLen)]))
	}
	return result
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

// private

func getBucketCounts(size, bucketSize int) int {
	if size%bucketSize == 0 {
		return size / bucketSize
	}
	return size/bucketSize + 1
}
