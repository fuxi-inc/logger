package logger

import (
	"os"
	"path"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// config.go contains all logger configs which supported by
// logrotate: https://gopkg.in/natefinch/lumberjack.v2
// zap: https://pkg.go.dev/go.uber.org/zap#hdr-Configuring_Zap

const (
	defaultFileName   = "default.log"
	defaultMaxSize    = 100000
	defaultMaxAge     = 7
	defaultMaxBackups = 4
	defaultRotate     = "1h"
	defaultLevel      = "info"
)

type conf struct {
	// configs for lumberjack
	FileName   string        `mapstructure:"log_file"`
	MaxSize    int           `mapstructure:"max_size"`
	MaxAge     int           `mapstructure:"max_age"`
	MaxBackups int           `mapstructure:"max_backups"`
	Rotate     string        `mapstructure:"rotate"` // rotate的间隔周期, 1h, 1d, 5min
	LevelStr   string        `mapstructure:"log_level"`
	Level      zapcore.Level `mapstructure:"-"`
}

func loadConfig(configPath string, mainKey string) (*conf, error) {
	if len(mainKey) == 0 {
		mainKey = "log"
	}

	c := new(conf)
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName("log") // 默认文件配置文件为app.toml

	if configPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		configPath = path.Join(dir, "log.toml")
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigFile(configPath)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	subv := v.Sub(mainKey)
	subv.SetDefault("log_file", defaultFileName)
	subv.SetDefault("max_size", defaultMaxSize)
	subv.SetDefault("max_age", defaultMaxAge)
	subv.SetDefault("max_backups", defaultMaxBackups)
	subv.SetDefault("rotate", defaultRotate)
	subv.SetDefault("log_level", defaultLevel)
	if err := subv.Unmarshal(c); err != nil {
		return nil, err
	}

	l, err := zapcore.ParseLevel(c.LevelStr)
	if err != nil {
		return nil, err
	}
	c.Level = l

	return c, nil
}
