package util

import (
	"time"

	"github.com/spf13/viper"
)

//variables
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	EmailServer         string        `mapstructure:"EMAIL_SERVER"`
	EmailServerPort     int           `mapstructure:"EMAIL_SERVER_PORT"`
	ClientUrl           string        `mapstructure:"CLIENT_URL"`
	EmailSenderAddress  string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
}

// LoadConfig reads config from file or env variables.

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
