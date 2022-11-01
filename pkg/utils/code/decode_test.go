package code

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnglishSecretEncodeDecode(t *testing.T) {
	secret := "this is a secret"
	encodeInt := EncodeSecret(secret)
	result := DecodeSecret(encodeInt)
	assert.Equal(t, secret, result)
}

func TestChineseSecretEncodeDecode(t *testing.T) {
	secret := "这是一个秘密"
	encodeInt := EncodeSecret(secret)
	result := DecodeSecret(encodeInt)
	assert.Equal(t, secret, result)
}

func TestKeyEncodeDecode(t *testing.T) {
	key := "ThisIsABigNumber"
	encodeInt, ok := EncodeKey(key)
	require.True(t, ok, "encode key to struct bit.Int failed")
	result := DecodeKey(encodeInt)
	assert.Equal(t, key, result)
}

func TestKeyEncodeDecodeKeys(t *testing.T) {
	key := "ThisIsABigNumber1_ThisIsABigNumber2"
	encodeIntegers, ok := EncodeKeys(key)
	require.True(t, ok, "encode key to keys failed")
	require.Equal(t, 2, len(encodeIntegers))
	result := DecodeKeys(encodeIntegers)
	assert.Equal(t, key, result)
}
