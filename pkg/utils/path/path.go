package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	KeyFilePrefix     = "shamir_"
	NecessaryFileName = KeyFilePrefix + "necessary-key"
	XKeyFilePrefix    = KeyFilePrefix + "x-key_"
	YKeyFilePrefix    = KeyFilePrefix + "y-key_"
)

// IsExist 返回路径是否存在
func IsExist(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsDir 返回路径是否是目录
func IsDir(path string) bool {
	stat, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// GetAllKeyFile 获取指定目录下所有的shamir key的文件路径
func GetAllKeyFile(path string) ([]string, error) {
	path = filepath.Clean(path)
	rd, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir fail: %w", err)
	}

	result := make([]string, 0)
	for _, fi := range rd {
		name := filepath.Base(fi.Name())
		if fi.IsDir() || !strings.HasPrefix(name, KeyFilePrefix) {
			continue
		}
		result = append(result, filepath.Join(path, name))
	}
	return result, nil
}

// IsNecessaryKeyExist 传入path是目录地址，判断该目录下是否存在necessary key
func IsNecessaryKeyExist(path string) bool {
	path = filepath.Join(filepath.Clean(path), NecessaryFileName)
	return IsExist(path)
}

func CheckNoKey(path string) error {
	files, err := GetAllKeyFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(files) != 0 {
		return fmt.Errorf("key files exist: %s", strings.Join(files, ","))
	}

	return nil
}

type KeyName struct {
	XKey string
	YKey string
}

// GetKeysName 从指定目录获取存在的密钥对的文件名，和必须密钥的文件名
func GetKeysName(path string) ([]*KeyName, string, error) {
	path = filepath.Clean(path)
	if !IsExist(path) {
		return nil, "", fmt.Errorf("path %q not exist", path)
	}

	if !IsNecessaryKeyExist(path) {
		return nil, "", fmt.Errorf("necessary key not exist")
	}

	names, err := GetAllKeyFile(path)
	if err != nil {
		return nil, "", err
	}
	namesMap := make(map[string]struct{}, len(names))
	for _, file := range names {
		namesMap[file] = struct{}{}
	}

	// 找到文件夹中的密钥对，密钥对的前缀分别是x和y相关前缀，后缀一致
	var keys []*KeyName
	for file := range namesMap {
		if !strings.HasPrefix(file, XKeyFilePrefix) {
			continue
		}

		suffix := strings.TrimPrefix(file, XKeyFilePrefix)
		if _, ok := namesMap[YKeyFilePrefix+suffix]; !ok {
			continue
		}
		keys = append(keys, &KeyName{
			XKey: XKeyFilePrefix + suffix,
			YKey: YKeyFilePrefix + suffix,
		})
	}

	if len(keys) == 0 {
		return nil, "", fmt.Errorf("path %q can not found key", path)
	}

	return keys, NecessaryFileName, nil
}
