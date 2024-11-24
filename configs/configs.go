package configs

import (
	"BTM-backend/pkg/osutil"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func init() {
	if err := os.Setenv("TZ", "UTC"); err != nil {
		panic(fmt.Errorf("fatal error config file: set time zone to utc: %w", err))
	}

	viper.AutomaticEnv()

	if _, isInCloudRun := os.LookupEnv("K_SERVICE"); !isInCloudRun {
		setFile()

		// Find and read the config file
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	loadConfig()
}

func setFile() {
	// 只有不在 cloud run 才能 load config
	path, err := osutil.FilePath("configs", "dev.env")
	if err != nil {
		viper.SetConfigName("app")
		viper.SetConfigType("env")
		// 因為執行位置在根目錄底下的話，就有可能是很多層
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		viper.AddConfigPath("../../configs")
		viper.AddConfigPath("../../../configs")
		viper.AddConfigPath("../../../../configs")
	} else {
		viper.SetConfigFile(path)
	}
}
