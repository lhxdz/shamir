package log

type Level string

const (
	InfoLevel   Level = "info"
	WarnLevel   Level = "warn"
	DebugLevel  Level = "debug"
	ErrorLevel  Level = "error"
	DPanicLevel Level = "dpanic"
	PanicLevel  Level = "panic"
)

var logger Logger

// SetGlobalLogger 传入 NewLogger 生成的 Logger 可以设置成全局logger
func SetGlobalLogger(sLogger Logger) {
	logger = sLogger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	if logger != nil {
		logger.Debug(args)
	}
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(args)
	}
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	if logger != nil {
		logger.Warn(args)
	}
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(args)
	}
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics
func DPanic(args ...interface{}) {
	if logger != nil {
		logger.DPanic(args)
	}
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	if logger != nil {
		logger.Panic(args)
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(template, args)
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	if logger != nil {
		logger.Infof(template, args)
	}
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	if logger != nil {
		logger.Warnf(template, args)
	}
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(template, args)
	}
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics
func DPanicf(template string, args ...interface{}) {
	if logger != nil {
		logger.DPanicf(template, args)
	}
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	if logger != nil {
		logger.Panicf(template, args)
	}
}
