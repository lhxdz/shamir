package code

import "math/big"

const (
	base = 56
)

// DecodeSecret 将解密后的秘密恢复成字符串
func DecodeSecret(secret *big.Int) string {
	if secret == nil {
		return ""
	}

	return string(secret.Bytes())
}

// DecodeKey 将加密生成的密钥输出成字符串
func DecodeKey(key *big.Int) string {
	if key == nil {
		return ""
	}

	return key.Text(base)
}
