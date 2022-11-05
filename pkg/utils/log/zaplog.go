package log

import (
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLever      = zapcore.InfoLevel
	defaultMaxSize    = 20
	defaultMaxBackups = 10
	defaultMaxAge     = 30
)

// 日志的输出位置
const (
	unDefine = iota
	console
	file
	both
)

var levelMap = map[Level]zapcore.Level{
	DebugLevel:  zapcore.DebugLevel,
	InfoLevel:   zapcore.InfoLevel,
	WarnLevel:   zapcore.WarnLevel,
	ErrorLevel:  zapcore.ErrorLevel,
	DPanicLevel: zapcore.DPanicLevel,
	PanicLevel:  zapcore.PanicLevel,
}

func getLogLevel(lvl Level) zapcore.Level {
	if level, ok := levelMap[Level(strings.ToLower(string(lvl)))]; ok {
		return level
	}
	return zapcore.InfoLevel
}

type Option func(conf *logConfig)

func WithLogLever(level Level) Option {
	return func(conf *logConfig) {
		conf.level = getLogLevel(level)
	}
}

// WithLogPath 输出到文件，可与 WithConsole 同时使用
func WithLogPath(path string) Option {
	return func(conf *logConfig) {
		conf.lumberjackConf.Filename = path
		if conf.logWriter == console {
			conf.logWriter = both
		} else if conf.logWriter != both {
			conf.logWriter = file
		}
	}
}

// WithConsole 输出到标准输出，可与 WithLogPath 同时使用
func WithConsole() Option {
	return func(conf *logConfig) {
		if conf.logWriter == file {
			conf.logWriter = both
		} else if conf.logWriter != both {
			conf.logWriter = console
		}
	}
}

// WithMaxSize 日志文件最大值，单位 M，必须使用 WithLogPath 才能生效
func WithMaxSize(size int) Option {
	return func(conf *logConfig) {
		conf.lumberjackConf.MaxSize = size
	}
}

// WithMaxBackups 日志备份文件最大个数，必须使用 WithLogPath 才能生效
func WithMaxBackups(maxBackups int) Option {
	return func(conf *logConfig) {
		conf.lumberjackConf.MaxBackups = maxBackups
	}
}

// WithMaxAge 日志备份文件保存的最大天数，必须使用 WithLogPath 才能生效
func WithMaxAge(maxAge int) Option {
	return func(conf *logConfig) {
		conf.lumberjackConf.MaxAge = maxAge
	}
}

func WithDefault() Option {
	return func(conf *logConfig) {
		conf.logWriter = unDefine
		conf.level = defaultLever
		conf.lumberjackConf.MaxSize = defaultMaxSize
		conf.lumberjackConf.MaxBackups = defaultMaxBackups
		conf.lumberjackConf.MaxAge = defaultMaxAge
		conf.lumberjackConf.Compress = true
	}
}

type logConfig struct {
	logWriter      int
	level          zapcore.Level
	lumberjackConf *lumberjack.Logger
}

func (l *logConfig) getZapCore() zapcore.Core {
	return zapcore.NewCore(getEncoder(), l.getWriteSyncer(), zap.NewAtomicLevelAt(l.level))
}

func (l *logConfig) getWriteSyncer() zapcore.WriteSyncer {
	syncers := make([]zapcore.WriteSyncer, 0)
	if l.logWriter != file {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	if l.logWriter >= file {
		syncers = append(syncers, zapcore.AddSync(l.lumberjackConf))
	}

	return zapcore.NewMultiWriteSyncer(syncers...)
}

func getEncoder() zapcore.Encoder {
	//获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func NewLogger(opts ...Option) Logger {
	conf := &logConfig{lumberjackConf: &lumberjack.Logger{}}
	WithDefault()(conf)
	for _, opt := range opts {
		opt(conf)
	}

	return zap.New(conf.getZapCore(), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
