package configs

import (
	"time"

	"github.com/spf13/viper"
)

var C Config

type Config struct {
	Port          string
	Gin           GinConfig
	TimeoutSecond time.Duration
}

type GinConfig struct {
	Mode string
}

func loadConfig() {
	C = Config{
		Port:          viper.GetString("port"),
		TimeoutSecond: viper.GetDuration("timeout_second"),
		Gin: GinConfig{
			Mode: viper.GetString("gin.mode"),
		},
	}
}
