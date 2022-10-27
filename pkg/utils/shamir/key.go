package shamir

import "math/big"

// Key 通过shamir方案加密后的密钥对 (x, y)，任意k个密钥对可以还原出密文
type Key struct {
	X *big.Int
	Y *big.Int
}
