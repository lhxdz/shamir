package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configInfo = `
# this is a config test
test:
  test1: 1
  test2: "test"
`

const (
	tmpDir = "/tmp"
)

func TestConfig(t *testing.T) {
	// 临时目录
	path, err := os.MkdirTemp(tmpDir, "")
	require.NoError(t, err)
	defer os.RemoveAll(path)
	// 修改config读取的路径
	configPath[0] = path
	// 临时config文件
	f, err := os.Create(filepath.Join(path, configFileName))
	require.NoError(t, err)
	_, err = f.Write([]byte(configInfo))
	require.NoError(t, err)

	err = InitConfig(false)
	require.NoError(t, err)

	assert.Equal(t, 1, viper.GetInt("test.test1"))
	assert.Equal(t, "test", viper.GetString("test.test2"))
}
