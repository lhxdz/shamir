package code

import (
	"bytes"
	"math/big"
)

const (
	base     = big.MaxBase
	splitKey = "_"
)

// DecodeSecret 将解密后的秘密恢复成字符串
func DecodeSecret(secret *big.Int) string {
	return string(getSecretBytes(secret))
}

func DecodeCompoundSecret(secret []*big.Int) string {
	if len(secret) == 0 {
		return ""
	}

	b := bytes.NewBuffer(make([]byte, 0, (len(secret[0].Bytes())-1)*len(secret)))
	for _, tmpSecret := range secret {
		b.Write(getSecretBytes(tmpSecret))
	}
	return b.String()
}

// DecodeKey 将加密生成的密钥输出成字符串
func DecodeKey(key *big.Int) string {
	if key == nil {
		return ""
	}

	return key.Text(base)
}

// DecodeKeys 将加密生成的密钥链输出成密钥字符串
func DecodeKeys(keys []*big.Int) string {
	result := ""
	for i, key := range keys {
		if i != 0 {
			result += splitKey
		}
		result += DecodeKey(key)
	}

	return result
}

// private

func getSecretBytes(secret *big.Int) []byte {
	if secret == nil || len(secret.Bytes()) < 1 {
		return []byte{}
	}

	// 去掉0xf前缀
	return secret.Bytes()[1:]
}
