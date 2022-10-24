package utils

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNextPrime(t *testing.T) {
	prime := big.NewInt(11)
	assert.Equal(t, int64(11), NextPrime(prime).Int64())
}
