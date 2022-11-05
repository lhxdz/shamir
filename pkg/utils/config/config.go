package config

import (
	"fmt"
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

func InitConfig(fast bool) error {
	if fast {
		return nil
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigName(configFileName)
	for _, path := range configPath {
		if !secure.ValidateFileSize(filepath.Join(path, configFileName), maxConfigFileSize) {
			fmt.Printf("Error: file%q size more than %dM", filepath.Join(path, configFileName), maxConfigFileSize/unitM)
			continue
		}
		viper.AddConfigPath(path)
	}
	viper.SetConfigType(defaultConfigType)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
