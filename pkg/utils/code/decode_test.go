package code

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnglishEncodeDecode(t *testing.T) {
	secret := "this is a secret"
	encodeInt := Encode(secret)
	result := Decode(encodeInt)
	assert.Equal(t, secret, result)
}

func TestChineseEncodeDecode(t *testing.T) {
	secret := "这是一个秘密"
	encodeInt := Encode(secret)
	result := Decode(encodeInt)
	assert.Equal(t, secret, result)
}
