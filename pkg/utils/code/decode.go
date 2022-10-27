package code

import "math/big"

func Decode(secret *big.Int) string {
	if secret == nil {
		return ""
	}

	return string(secret.Bytes())
}
