package math

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPrime(t *testing.T) {
	prime := big.NewInt(10)
	assert.Equal(t, int64(11), NextPrime(prime).Int64())
}
func TestFastPrime(t *testing.T) {
	prime := fastPrimes[len(fastPrimes)-2]
	assert.Equal(t, 0, FastPrime(prime).Cmp(fastPrimes[len(fastPrimes)-1]))
}
