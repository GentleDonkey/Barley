package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

var JwtKey = []byte("my_secret_key")

const (
	defaultDbDriver      = "MySQL"
	defaultDbSource      = "root:19900718qzyQZY@tcp(localhost:3306)/barley?charset=utf8&parseTime=True&loc=Local"
	defaultServerAddress = "localhost:8080"
	defaultVersion       = "/api/v1"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Version       string `mapstructure:"VERSION"`
}

var C Config

func LoadConfig() Config {
	v := viper.New()
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()
	v.SetDefault("SERVER_ADDRESS", defaultServerAddress)
	v.SetDefault("VERSION", defaultVersion)
	v.SetDefault("DB_SOURCE", defaultDbSource)
	v.SetDefault("DB_DRIVER", defaultDbDriver)
	err := v.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("fatal error on init config"))
	}
	return C
}
