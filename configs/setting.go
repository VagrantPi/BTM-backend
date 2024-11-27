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
	Db            DbConfig
	JWT           JWTConfig
}

type GinConfig struct {
	Mode string
}

type JWTConfig struct {
	Secret string
}

type DbConfig struct {
	Debug              bool
	Dialect            string
	Host               string
	Port               string
	Username           string
	Password           string
	Flag               string
	SslMode            string
	Database           string
	ConnMaxLifetimeSec int
	MaxOpenConns       int
	MaxIdleConns       int
}

func loadConfig() {
	C = Config{
		Port:          viper.GetString("port"),
		TimeoutSecond: viper.GetDuration("timeout_second"),
		Gin: GinConfig{
			Mode: viper.GetString("gin.mode"),
		},
		Db: DbConfig{
			Debug:              viper.GetBool("db.debug"),
			Dialect:            viper.GetString("db.dialect"),
			Host:               viper.GetString("db.host"),
			Port:               viper.GetString("db.port"),
			Username:           viper.GetString("db.username"),
			Password:           viper.GetString("db.password"),
			Flag:               viper.GetString("db.flag"),
			SslMode:            viper.GetString("db.ssl_mode"),
			Database:           viper.GetString("db.database"),
			MaxOpenConns:       viper.GetInt("db.max_open_conns"),
			MaxIdleConns:       viper.GetInt("db.max_idle_conns"),
			ConnMaxLifetimeSec: viper.GetInt("db.conn_max_life_time_sec"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("jwt.secret"),
		},
	}
}
