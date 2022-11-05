package config

import (
	"path/filepath"

	"github.com/spf13/viper"

	"shamir/pkg/utils/secure"
)

const (
	unitM = 1024 * 1024
	// 默认最大配置文件大小 20M
	maxConfigFileSize = 20 * unitM

	configFileName = "shamir.yaml"

	envPrefix = "SHAMIR"

	defaultConfigType = "yaml"
)

var configPath = []string{"/var/lib/shamir", "/usr/lib/shamir"}

// InitConfig 当传入参数fast==true时，将使用默认的config
func InitConfig(fast bool) error {
	defer SetDefault()
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	if fast {
		return nil
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType(defaultConfigType)
	for _, path := range configPath {
		if !secure.ValidateFileSize(filepath.Join(path, configFileName), maxConfigFileSize) {
			continue
		}

		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 配置文件未找到错误，暂时忽略
			return err
		}
	}

	return nil
}
