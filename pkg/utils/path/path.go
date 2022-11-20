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
		fmt.Println("read dir fail:", err)
		return nil, err
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
