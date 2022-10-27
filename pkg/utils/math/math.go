package math

import "math"

// InvMod 计算 a^(-1) mod n. While panic if undefined.
func InvMod(a, n int64) int64 {

	d, s := XGCD(a, n)
	if d != 1 {
		panic(any("InvMod: inverse undefined"))
	}
	if s < 0 {
		return s + n
	} else {
		return s
	}
}

func XGCD(a, b int64) (int64, int64) {
	if a < -math.MaxInt64 || b < -math.MaxInt64 {
		panic(any("XGCD: integer overflow"))
	}

	var aNeg int64 = 1
	if a < 0 {
		a = -a
		aNeg = -1
	}

	if b < 0 {
		b = -b
	}

	var u1, v1, u2, v2 int64 = 1, 0, 0, 1
	for b != 0 {
		q, r := a/b, a%b
		a, b = b, r
		u0, v0 := u2, v2
		u2, v2 = u1-q*u2, v1-q*v2
		u1, v1 = u0, v0
	}

	return a, u1 * aNeg
}
