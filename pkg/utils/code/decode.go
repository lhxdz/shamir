package code

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"math/big"
)

const (
	base     = big.MaxBase
	splitKey = "_"
)

// DecodeSecret 将解密后的秘密恢复成字符串
func DecodeSecret(secret *big.Int) string {
	return string(getSecretBytes(secret))
}

func DecodeCompoundSecret(secret []*big.Int) string {
	if len(secret) == 0 {
		return ""
	}

	b := bytes.NewBuffer(make([]byte, 0, (len(secret[0].Bytes())-1)*len(secret)))
	for _, tmpSecret := range secret {
		b.Write(getSecretBytes(tmpSecret))
	}
	return b.String()
}

// DecodeKey 将加密生成的密钥输出成字符串
func DecodeKey(key *big.Int) string {
	if key == nil {
		return ""
	}

	return key.Text(base)
}

// DecodeKeys 将加密生成的密钥链输出成密钥字符串
func DecodeKeys(keys []*big.Int) string {
	result := ""
	for i, key := range keys {
		if i != 0 {
			result += splitKey
		}
		result += DecodeKey(key)
	}

	return result
}

type SecretDecoder struct {
	// 用于记录最后一次传入的数据，因为最后一个是hash校验值
	lastData []byte
	hash     hash.Hash
	writer   io.Writer
}

func NewSecretDecoder(writer io.Writer) *SecretDecoder {
	hashWriter := sha256.New()
	multiWriter := io.MultiWriter(writer, hashWriter)
	return &SecretDecoder{
		writer: multiWriter,
		hash:   hashWriter,
	}
}

func (s *SecretDecoder) Write(data *big.Int) error {
	if data == nil || len(data.Bytes()) == 0 {
		return fmt.Errorf("invalid data bit integer")
	}

	// 每次写上次的数据，防止将最后一个hash校验的数据误写入
	if s.lastData != nil {
		n, err := s.writer.Write(s.lastData)
		if err != nil {
			return fmt.Errorf("write data failed: %w", err)
		}
		if n != len(s.lastData) {
			return fmt.Errorf("write data failed, expected write %d bytes, actual %d bytes", len(s.lastData), n)
		}
	}

	// 将这次传入的数据保存
	s.lastData = getSecretBytes(data)
	return nil
}

func (s *SecretDecoder) HashCheck() error {
	expectedHash := string(s.lastData)
	actuallyHash := string(s.hash.Sum(nil))
	if expectedHash != actuallyHash {
		return fmt.Errorf("expected hash %s, actual hash %s", expectedHash, actuallyHash)
	}

	return nil
}

type KeyDecoder struct {
	split  string
	writer io.Writer
}

func NewKeyDecoder(writer io.Writer) *KeyDecoder {
	return &KeyDecoder{
		writer: writer,
	}
}

func (k *KeyDecoder) Write(key *big.Int) error {
	if key == nil {
		return fmt.Errorf("%w, nil point", InvalidKey)
	}

	// 首次写入时前面没有分隔符
	keyData := key.Append([]byte(k.split), base)
	n, err := k.writer.Write(keyData)
	if err != nil {
		return fmt.Errorf("write key data failed: %w", err)
	}
	if n != len(keyData) {
		return fmt.Errorf("write key data failed, expected write %d bytes, actual %d bytes", len(keyData), n)
	}

	// 首次写入后，添加分隔符
	k.split = splitKey
	return nil
}

// private

func getSecretBytes(secret *big.Int) []byte {
	if secret == nil || len(secret.Bytes()) < 1 {
		return []byte{}
	}

	// 去掉0xf前缀
	return secret.Bytes()[1:]
}
