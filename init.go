package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	encoderConfig = zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "time",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		CallerKey:      "location",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder, EncodeCaller: zapcore.ShortCallerEncoder,
	}
)

// MustInit a logger instance, otherwise panic err msgs
func MustInit(logConfigPath string) Logger {
	l, err := Init(logConfigPath)

	if err != nil {
		panic(err.Error())
	}

	return l
}

// Init 初始化配置日志模块
func Init(logConfigPath string) (Logger, error) {
	// 首先使用conf加载log配置
	ret := new(zapL)
	c, err := loadConfig(logConfigPath, "log")
	if err != nil {
		return nil, err
	}
	if s, err := loadServiceLogger(c); err != nil {
		return nil, err
	} else {
		ret.serviceLogger = s
	}

	c, err = loadConfig(logConfigPath, "errlog")
	if err != nil {
		return nil, err
	}
	if s, err := loadServiceLogger(c); err != nil {
		return nil, err
	} else {
		ret.serviceErrLogger = s
	}

	return ret, nil
}

func loadServiceLogger(c *conf) (*zap.Logger, error) {
	lj := &lumberjack.Logger{
		Filename:   c.FileName,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		LocalTime:  true,
	}

	// 按照c.Rotate启动一个time ticker，然后定期执行Rotate
	duration, err := time.ParseDuration(c.Rotate)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(duration)

	go func(lj *lumberjack.Logger) {
		defer ticker.Stop()
		for range ticker.C {
			if err := rotate(lj); err != nil {
				panic(err.Error())
			}
		}
	}(lj)

	writer := zapcore.AddSync(lj)

	// todo，后续支持可以配置多个core
	var cores = []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), c.Level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(zapcore.AddSync(os.Stdout)), c.Level),
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1)), nil
}

var rotate = func(lj *lumberjack.Logger) error {
	if err := lj.Rotate(); err != nil {
		return err
	}

	return nil
}
