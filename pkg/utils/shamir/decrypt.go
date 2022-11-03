package shamir

import "math/big"

const (
	negative = -1
)

// Decrypt 根据密钥对解密，传入 len(keys) 必须是门限值threshold, prime 也必须是加密时使用的素数。
// keys要求不能有重复元素，否则panic
func Decrypt(keys []Key, prime *big.Int) *big.Int {
	result := big.NewInt(0)
	for i, key := range keys {
		tmp := product(getXKeysExceptI(keys, i), key.X, prime)
		tmp = tmp.Mul(tmp, key.Y)
		result = result.Add(result, tmp)
		result = result.Mod(result, prime)
	}
	return result
}

func getXKeysExceptI(keys []Key, i int) []*big.Int {
	xKeys := make([]*big.Int, 0, len(keys)-1)
	for j, key := range keys {
		if j == i {
			continue
		}
		xKeys = append(xKeys, key.X)
	}

	return xKeys
}

// 求((-xKeys[0])*...(-xKeys[n])) / ((xKeysI - xKeys[0])*...(xKeysI - xKeys[n])) mod prime
// 其中 n = (len(xKeys)-1)
// 其中的xKey不能为0
// prime为0会panic
// 若其中xKeys某一项和xKeysI值相等，则会panic
func product(xKeys []*big.Int, xKeysI *big.Int, prime *big.Int) *big.Int {
	result := big.NewInt(1)
	if len(xKeys)%2 != 0 {
		// 奇数个负数分子相乘为奇数
		result = big.NewInt(-1)
	}
	for _, key := range xKeys {
		result = result.Mul(result, key)
		result = result.Mod(result, prime)
		denominator := new(big.Int).Sub(xKeysI, key)
		if denominator.Sign() == negative {
			// 若分母为负数,则应先求模，再求 x**y mod p
			denominator = denominator.Mod(denominator, prime)
		}
		denominator = denominator.Exp(denominator, big.NewInt(-1), prime)
		// 分子分母取模结果相乘再取模
		result = result.Mul(result, denominator)
		result = result.Mod(result, prime)
	}

	return result
}
