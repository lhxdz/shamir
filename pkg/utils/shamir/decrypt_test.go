package shamir

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type decryptEncryptSuit struct {
	suite.Suite

	secret                *big.Int
	keys                  []Key
	threshold, keysNumber int
	prime                 *big.Int
}

func (d *decryptEncryptSuit) SetupSuite() {
	d.secret, _ = new(big.Int).SetString("TPPtg4XDmRfUza5cLAvfE2UQJenIvUBietpy4QNftZGaoQY2S1UAUKcO02EwpTmRK06ES82cTuok7046TKORxyztNBaCA0558u5ErFKZ9PKv1EGaeYcAlISTfHvevHFkURIzvfRQZ0SzkO7YJ4xjzxMQedR0EyPk6XPduHkUtQaP0n0n0mPzeUXNYoZWzQiTBdYffGpEbEcOeBhoEJqPc82HHlKsWIqXpN1NiuyH4QwhaGIzkWfUls0FeO3u8mLQZ95pEsz3wehziixypvXfkUkXasGM0B6OeHcIWjnc5sQFLLBC72UWIZK1ZokDKsPX5kMlF4PsczJbxTwTA5PScU6fL3nw3bZ8qUmMRGQTo8G1lYZhSu9bWyPhlv6Xs0nEJzSj5miOhvVCKpLXws0yyCkuJJIzYc8nHXZUzxosdMur619pIQCkO9c6O5dDUw13YqtVRTkprQBXjufd0uTk7UjE6avSijlaJI5C6R2SNV58HpVmvQa680QAPByzkxVfKN1Rp0zHKg30cZ5aadq9w0Scmttv6Q2DRV1IB8HXK3kGcPhI6GRJi23Q2RISetEdGcuI4YKGXUeYQZNS1YQoiQa54S2N2Xszt1ypQQ37rbkHoDBO5nJlXH2BoKl0JDW8EXhYZ59RNUzgNjFnzJ46pognVWo7JO97iClK9XRvAyejzERe9Krw0YNmm2hg6Ov5YflBaXYRClvYcUi2quiNtjEsskvGCVrMMSL", big.MaxBase)
	d.threshold, d.keysNumber = 4, 10

	keys, prime, err := Encrypt(d.secret, d.threshold, d.keysNumber, true)
	if err != nil {
		d.FailNowf("get encrypt of secret(%s) failed: %d", d.secret.Text(big.MaxBase), err)
	}

	d.keys, d.prime = keys, prime
}

func (d *decryptEncryptSuit) TestDecrypt() {
	result := Decrypt(d.keys[:d.threshold], d.prime)
	assert.Equal(d.T(), 0, result.Cmp(d.secret))
}

func TestShamir(t *testing.T) {
	test := new(decryptEncryptSuit)
	suite.Run(t, test)
}
