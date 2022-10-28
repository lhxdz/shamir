package code

import "math/big"

// EncodeSecret 将字符串秘密编码成为大整数，方便加密
func EncodeSecret(secret string) *big.Int {
	return new(big.Int).SetBytes([]byte(secret))
}

// EncodeKey 将密钥字符串恢复成大整数
func EncodeKey(key string) (*big.Int, bool) {
	return new(big.Int).SetString(key, base)
}
