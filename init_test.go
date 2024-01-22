package logger

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func printrotate(lj *lumberjack.Logger) error {
	fmt.Printf("timer of %s, called at %s \n", lj.Filename, time.Now().String())
	return nil
}

func Test_loadServiceLogger(t *testing.T) {
	rotate = printrotate
	type args struct {
		c *conf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"normal_5s",
			args{
				c: &conf{
					FileName:   "normal_5s",
					MaxSize:    1000,
					MaxAge:     10,
					MaxBackups: 3,
					Rotate:     "5s",
					Level:      zapcore.InfoLevel,
				},
			},
			false,
		},

		{
			"error_10s",
			args{
				c: &conf{
					FileName:   "error_10s",
					MaxSize:    1000,
					MaxAge:     10,
					MaxBackups: 3,
					Rotate:     "10s",
					Level:      zapcore.ErrorLevel,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadServiceLogger(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadServiceLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	time.Sleep(time.Minute)
}
