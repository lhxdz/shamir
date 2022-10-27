package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
)

type decryptEncryptSuit struct {
	suite.Suite

	secret                *big.Int
	keys                  []Key
	threshold, keysNumber int
	prime                 *big.Int
}

func (d *decryptEncryptSuit) SetupSuite() {
	d.secret = big.NewInt(123456789)
	d.threshold, d.keysNumber = 3, 10

	keys, prime, err := Encrypt(d.secret, d.threshold, d.keysNumber)
	if err != nil {
		d.FailNowf("get encrypt of secret(%s) failed: %d", d.secret.String(), err)
	}

	d.keys, d.prime = keys, prime
}

func (d *decryptEncryptSuit) TestDecrypt() {
	result := Decrypt(d.keys[:d.threshold], d.prime)
	assert.Equal(d.T(), result.Cmp(d.secret), 0)
}

func TestShamir(t *testing.T) {
	test := new(decryptEncryptSuit)
	suite.Run(t, test)
}
