package code

import "math/big"

// Encode 将字符串秘密编码成为大整数，方便加密
func Encode(secret string) *big.Int {
	return new(big.Int).SetBytes([]byte(secret))
}
