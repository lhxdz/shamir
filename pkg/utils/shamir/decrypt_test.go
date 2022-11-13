package shamir

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/compute"
)

const (
	bigSecret = `
shamir-tools
this is a tool for shamir (k, n) threshold scheme

使用方式
加密：
输入门限值t、密钥个数n、秘密 即可加密，将会生成一个必须密钥necessary_key、n个密钥(每个密钥包含x、y)，其中任意t个密钥可以恢复秘密

lhx@DESKTOP-0GALLEM:~$ shamir encrypt -t 2 -n 3 "this is a secret.同时可以使用中文。"
necessary key: 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_MV0yI9PnICUSbyxh9iceG0W0HlBn6PDFcfsDDWIHFs3
+-------------------------+-----------------------------------------------------------------------------------------------------------+
|          KEY X          |                                                   KEY Y                                                   |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| pBE2KvVaYN_569FAGOeEKr  | BVBszEaZArdp0XuVWgxQkDUSr6ylDCgi5t2i7nH4nMEh9ehZeyrPSbs8OnOu_pENXjVJguH45QnxV4k6P5gkxfhm989b23uMqVCd5dWc  |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| 8PXwCwkVb7L_2huJ6o8FBgE | uvzHiQfto4zzEJvjhKhCd30XIosCKe2XR3Uh8ULKXULqgA9skmQqlQ0yTWHk_vf0gf0Elij0iedGYRKHrbpTg5CnYvItrt1J9NalREqy  |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| aWTt0gYv09K_8Id7tS4dstb | 1CvtrfulBa0ObWofpE4ygNK4bQdDWetizrXSXAdWbLTC4T92YAlJpLvFvkMdx_iQcV2Pdrfgb6iCbv9zoBDpXfekqW1f5fVzj7gs3WSCp |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
解密：
分别输入t个密钥x和y，其中第i个x和第i个y是一对密钥，同时输入necessary_key，即可恢复秘密

lhx@DESKTOP-0GALLEM:~$ shamir decrypt -n 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_MV0yI9PnICUSbyxh9iceG0W0HlBn6PDFcfsDDWIHFs3 -x pBE2KvVaYN_569FAGOeEKr,8PXwCwkVb7L_2huJ6o8FBgE -y BVBszEaZArdp0XuVWgxQkDUSr6ylDCgi5t2i7nH4nMEh9ehZeyrPSbs8OnOu_pENXjVJguH45QnxV4k6P5gkxfhm989b23uMqVCd5dWc,uvzHiQfto4zzEJvjhKhCd30XIosCKe2XR3Uh8ULKXULqgA9skmQqlQ0yTWHk_vf0gf0Elij0iedGYRKHrbpTg5CnYvItrt1J9NalREqy
this is a secret.同时可以使用中文。
或者：

lhx@DESKTOP-0GALLEM:~$ shamir decrypt -x 1ARgxTZ9eRv -x 3uo9qE5tjus -y a7BG2cYFbHFGbcSYe3g90lAyJM6HneKZJqi4U0sY6WAWbUtxVhNbFk -y cWDRLeRNO7swCBis1MkEpWNCC4J4Mbz
my2NvEMiVCkO209SP7DesrI -n lwegNercKn81SX8boBDos0ksm0cQpaHm7Nu6qJPmJMk2Kgz25o4ggX
this is a secret\!中文也可以加密
详细介绍：
背景
现在大家很熟悉的加解密场景，可以看做A hash成 B，由B可以还原成原来的A

但是在某些场景下，需要将密钥拆分成许多子密钥，让每一方分别持有一个子密钥；需要密钥时，必须有一定数量的子密钥合在一起才能恢复出整个密钥。这样的过程称为密钥共享，更广义的，也可称为秘密共享

使用秘密共享方案进行密钥管理，有如下优点：

解决主密钥丢失即所有密钥丢失的问题
避免出现权限最高的个体
提高了密钥的安全性，因为攻击者必须获得足够多的子密钥才能恢复出主密钥
所有秘密共享方案中，最经典的要数Shamir于1979年提出的 (k,n)门限秘密共享方案

方案简述：
将秘密S分成n个子秘密
任何 k 个子秘密合在一起能恢复出秘密S
任何 <k 个子秘密合在一起无法恢复出秘密S
`
)

type decryptEncryptSuit struct {
	suite.Suite

	threshold, keysNumber int
	chunkSize             int

	secret *big.Int
	keys   []code.Key
	prime  *big.Int

	bigSecret []*big.Int
	bigKeys   []code.CompoundKey
	bigPrime  []*big.Int
}

func (d *decryptEncryptSuit) SetupSuite() {
	// 小秘密
	d.secret, _ = new(big.Int).SetString("TPPtg4XDmRfUza5cLAvfE2UQJenIvUBietpy4QNftZGaoQY2S1UAUKcO02EwpTmRK06ES82cTuok7046TKORxyztNBaCA0558u5ErFKZ9PKv1EGaeYcAlISTfHvevHFkURIzvfRQZ0SzkO7YJ4xjzxMQedR0EyPk6XPduHkUtQaP0n0n0mPzeUXNYoZWzQiTBdYffGpEbEcOeBhoEJqPc82HHlKsWIqXpN1NiuyH4QwhaGIzkWfUls0FeO3u8mLQZ95pEsz3wehziixypvXfkUkXasGM0B6OeHcIWjnc5sQFLLBC72UWIZK1ZokDKsPX5kMlF4PsczJbxTwTA5PScU6fL3nw3bZ8qUmMRGQTo8G1lYZhSu9bWyPhlv6Xs0nEJzSj5miOhvVCKpLXws0yyCkuJJIzYc8nHXZUzxosdMur619pIQCkO9c6O5dDUw13YqtVRTkprQBXjufd0uTk7UjE6avSijlaJI5C6R2SNV58HpVmvQa680QAPByzkxVfKN1Rp0zHKg30cZ5aadq9w0Scmttv6Q2DRV1IB8HXK3kGcPhI6GRJi23Q2RISetEdGcuI4YKGXUeYQZNS1YQoiQa54S2N2Xszt1ypQQ37rbkHoDBO5nJlXH2BoKl0JDW8EXhYZ59RNUzgNjFnzJ46pognVWo7JO97iClK9XRvAyejzERe9Krw0YNmm2hg6Ov5YflBaXYRClvYcUi2quiNtjEsskvGCVrMMSL", big.MaxBase)
	d.threshold, d.keysNumber = 4, 10

	keys, prime, err := Encrypt(d.secret, d.threshold, d.keysNumber, true)
	require.NoErrorf(d.T(), err, "get encrypt of secret(%s) failed", d.secret.Text(big.MaxBase))
	d.keys, d.prime = keys, prime

	// 大型秘密
	d.chunkSize = 100
	for i := 0; i < len(bigSecret); i += d.chunkSize {
		d.bigSecret = append(d.bigSecret, code.EncodeSecret(bigSecret[i:compute.Min(i+d.chunkSize, len(bigSecret))]))
	}
	bigKeys, bigPrime, err := HashEncrypt(d.bigSecret, d.threshold, d.keysNumber, true)
	require.NoError(d.T(), err, "get hash encrypt of big secret failed")
	d.bigKeys, d.bigPrime = bigKeys, bigPrime
}

func (d *decryptEncryptSuit) TestDecrypt() {
	// 小秘密
	result, err := Decrypt(d.keys[:d.threshold], d.prime)
	assert.NoError(d.T(), err)
	assert.Equal(d.T(), 0, result.Cmp(d.secret))

	// 大型秘密
	bigResultInt, err := HashDecrypt(d.bigKeys[:d.threshold], d.bigPrime)
	assert.NoError(d.T(), err)
	bigResult := ""
	for _, tmp := range bigResultInt {
		bigResult += code.DecodeSecret(tmp)
	}
	assert.Equal(d.T(), bigSecret, bigResult)
}

func TestShamir(t *testing.T) {
	test := new(decryptEncryptSuit)
	suite.Run(t, test)
}
