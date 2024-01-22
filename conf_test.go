package logger

import (
	"os"
	"testing"
)

var (
	confData = `
	[log]
	log_file="./log/dip.log"
	log_level="Debug"
	rotate = "5s"
	max_size = 100
	max_age = 3
	max_backups = 3
	`
)

func Test_loadConfig(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "conf")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = f.Write([]byte(confData))
	if err != nil {
		t.Fatal(err.Error())
	}
	f.Close()

	type args struct {
		configPath string
		mainKey    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"default_case",
			args{
				f.Name(),
				"log",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadConfig(tt.args.configPath, tt.args.mainKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
