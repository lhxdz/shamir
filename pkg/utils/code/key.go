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

type CompoundKey struct {
	X []*big.Int
	Y []*big.Int
}

type StrKey struct {
	X string `json:"key_x" yaml:"key_x"`
	Y string `json:"key_y" yaml:"key_y"`
}

func EncodeAbleKey(key Key) *StrKey {
	return &StrKey{
		X: DecodeKey(key.X),
		Y: DecodeKey(key.Y),
	}
}

func EncodeAbleCompoundKeys(keys []CompoundKey) []*StrKey {
	result := make([]*StrKey, 0, len(keys))
	for _, key := range keys {
		result = append(result, &StrKey{
			X: DecodeKeys(key.X),
			Y: DecodeKeys(key.Y),
		})
	}
	return result
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
			return nil, fmt.Errorf("invalid x key: %q", key.X)
		}
		yKey, ok := EncodeKey(key.Y)
		if !ok {
			return nil, fmt.Errorf("invalid y key: %q", key.Y)
		}
		result = append(result, Key{X: xKey, Y: yKey})
	}
	return result, nil
}

func EncodeStrCompoundKeys(keys []*StrKey) ([]CompoundKey, error) {
	result := make([]CompoundKey, 0, len(keys))
	for _, key := range keys {
		xKey, ok := EncodeKeys(key.X)
		if !ok {
			return nil, fmt.Errorf("invalid x key: %q", key.X)
		}
		yKey, ok := EncodeKeys(key.Y)
		if !ok {
			return nil, fmt.Errorf("invalid y key: %q", key.Y)
		}
		result = append(result, CompoundKey{X: xKey, Y: yKey})
	}
	return result, nil
}
