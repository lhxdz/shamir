package code

import (
	"fmt"
	"math/big"
)

// Key 通过shamir方案加密后的密钥对 (x, y)，任意k个密钥对可以还原出密文
type Key struct {
	X *big.Int
	Y *big.Int
}

type StrKey struct {
	X string
	Y string
}

func EncodeAbleKey(key Key) *StrKey {
	return &StrKey{
		X: DecodeKey(key.X),
		Y: DecodeKey(key.Y),
	}
}

func EncodeAbleKeys(keys []Key) []*StrKey {
	result := make([]*StrKey, 0, len(keys))
	for _, key := range keys {
		result = append(result, &StrKey{
			X: DecodeKey(key.X),
			Y: DecodeKey(key.Y),
		})
	}
	return result
}

func EncodeStrKeys(keys []*StrKey) ([]Key, error) {
	result := make([]Key, 0, len(keys))
	for _, key := range keys {
		xKey, ok := EncodeKey(key.X)
		if !ok {
			return nil, fmt.Errorf("invalid x key: %s", key.X)
		}
		yKey, ok := EncodeKey(key.Y)
		if !ok {
			return nil, fmt.Errorf("invalid y key: %s", key.Y)
		}
		result = append(result, Key{X: xKey, Y: yKey})
	}
	return result, nil
}
