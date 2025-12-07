package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var Log *zap.Logger

func Init() {
	runMode := "debug"
	if runMode == "debug" {
		core := zapcore.NewTee(
			zapcore.NewCore(getEncoder(), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
		Log = zap.New(core, zap.AddCaller())
	} else {
		fileLog()
	}
}
func fileLog() {
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.InfoLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.WarnLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.ErrorLevel
	})
	panicPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.PanicLevel
	})
	fatalPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zapcore.FatalLevel
	})
	cores := [...]zapcore.Core{
		getEncoderCore("debug.log", debugPriority),
		getEncoderCore("info.log", infoPriority),
		getEncoderCore("warn.log", warnPriority),
		getEncoderCore("error.log", errorPriority),
		getEncoderCore("panic.log", panicPriority),
		getEncoderCore("fatal.log", fatalPriority),
	}
	Log = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
}
func getEncoderCore(fileName string, level zapcore.LevelEnabler) zapcore.Core {
	writer := getLogWriter(fileName)
	return zapcore.NewCore(getEncoder(), writer, level)
}
func getLogWriter(fileName string) zapcore.WriteSyncer {
	logDir := "./runtime/logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		println("create log error" + err.Error())
		return zapcore.AddSync(os.Stdout)
	}
	logPath := path.Join(logDir, fileName)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 设置日志格式
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(getEncoderConfig())
}
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
