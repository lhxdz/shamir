package code

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strings"

	"github.com/pkg/errors"

	"shamir/pkg/utils/compute"
)

const (
	bytePrefix = 0xff
	// 目前最大分割的子秘密长1024，其生成的key不会超过2048长度
	maxKeyLen = 2048
)

var (
	InvalidKey = errors.New("invalid key")
)

// EncodeSecret 将字符串秘密编码成为大整数，方便加密
func EncodeSecret(secret string) *big.Int {
	// 所有的秘密都加上 0xff前缀, 避免全零数据丢失真实数据
	return new(big.Int).SetBytes(append([]byte{bytePrefix}, secret...))
}

func EncodeCompoundSecret(secret string, splitLen int) []*big.Int {
	result := make([]*big.Int, 0, getBucketCounts(len(secret), splitLen))
	for i := 0; i < len(secret); i += splitLen {
		result = append(result, EncodeSecret(secret[i:compute.Min(len(secret), i+splitLen)]))
	}
	return result
}

// EncodeKey 将密钥字符串恢复成大整数密钥
func EncodeKey(key string) (*big.Int, bool) {
	return new(big.Int).SetString(key, base)
}

// EncodeKeys 将密钥字符串恢复成大整数密钥链
func EncodeKeys(key string) ([]*big.Int, bool) {
	subKeys := strings.Split(key, splitKey)
	bigIntegers := make([]*big.Int, 0, len(subKeys))
	for _, subKey := range subKeys {
		subInteger, ok := EncodeKey(subKey)
		if !ok {
			return nil, false
		}
		bigIntegers = append(bigIntegers, subInteger)
	}
	return bigIntegers, true
}

type SecretEncoder struct {
	splitLen int
	hash     hash.Hash
	reader   io.Reader
}

func NewSecretEncoder(reader io.Reader, splitLen int) *SecretEncoder {
	return &SecretEncoder{
		splitLen: splitLen,
		hash:     sha256.New(),
		reader:   reader,
	}
}

func (s *SecretEncoder) Read() (*big.Int, error) {
	data := make([]byte, s.splitLen)
	n, err := s.reader.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("read secret file failed: %w", err)
	}

	if n == 0 {
		return nil, nil
	}

	nHash, err := s.hash.Write(data[:n])
	if err != nil {
		return nil, fmt.Errorf("secret hash check failed: %w", err)
	}
	if nHash != n {
		return nil, fmt.Errorf("secret hash check failed, hash write expected %d bytes, actual %d bytes", n, nHash)
	}

	return new(big.Int).SetBytes(append([]byte{bytePrefix}, data[:n]...)), nil
}

func (s *SecretEncoder) GetHash() *big.Int {
	return new(big.Int).SetBytes(s.hash.Sum(nil))
}

type KeyEncoder struct {
	reader *bufio.Reader
}

func NewKeyEncoder(reader io.Reader) *KeyEncoder {
	return &KeyEncoder{
		reader: bufio.NewReaderSize(reader, maxKeyLen),
	}
}

// Read 返回密钥，密钥类型(是否是hash值的密钥)
func (s *KeyEncoder) Read() (*big.Int, bool, error) {
	data, err := s.reader.ReadSlice(splitKey[0])
	if err != nil && !errors.Is(err, io.EOF) {
		// 正确的key的长度无法达到这么大
		if errors.Is(err, bufio.ErrBufferFull) {
			err = InvalidKey
		}
		return nil, false, fmt.Errorf("read key file failed: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.Is(err, io.EOF), InvalidKey
	}

	// 非hash值密钥
	if !errors.Is(err, io.EOF) {
		// 去除最后一个字符的分隔符
		data = data[:len(data)-1]
	}

	result, e := encodeKey(data)
	if e != nil {
		return nil, errors.Is(err, io.EOF), e
	}

	return result, errors.Is(err, io.EOF), nil
}

// private

func getBucketCounts(size, bucketSize int) int {
	if size%bucketSize == 0 {
		return size / bucketSize
	}
	return size/bucketSize + 1
}

func encodeKey(data []byte) (*big.Int, error) {
	if len(data) == 0 {
		return nil, InvalidKey
	}

	result, ok := EncodeKey(string(data))
	if !ok {
		return nil, InvalidKey
	}
	return result, nil
}
