package configure

import "github.com/spf13/viper"

type Config struct {
	DBUrl             string `mapstructure:"DATABASE_URL"`
	VideoDir          string `mapstructure:"VIDEO_DIR"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
