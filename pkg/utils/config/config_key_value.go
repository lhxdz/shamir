package config

import "github.com/spf13/viper"

const (
	// log相关配置

	LogPathKey     = "log.path"
	defaultLogPath = "/var/log/shamir/shamir.log"
	// LogMaxSize 日志文件最大大小
	LogMaxSize     = "log.max_size"
	defaultMaxSize = 30 // 单位M
	// LogMaxBackups 日志文件最大备份个数
	LogMaxBackups     = "log.max_backups"
	defaultMaxBackups = 10
	// LogMaxAge 日志备份文件最大存在天数
	LogMaxAge     = "log.max_age"
	defaultMaxAge = 60
)

func SetDefault() {
	viper.SetDefault(LogPathKey, defaultLogPath)
	viper.SetDefault(LogMaxSize, defaultMaxSize)
	viper.SetDefault(LogMaxBackups, defaultMaxBackups)
	viper.SetDefault(LogMaxAge, defaultMaxAge)
}
