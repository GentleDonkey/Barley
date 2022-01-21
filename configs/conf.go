package configs

import (
	"github.com/spf13/viper"
)

var JwtKey = []byte("my_secret_key")

//var Host = "localhost:"
//var RestPath = "8080"
//var Version = "/api/v1"

//var DbPath = "3306"
//var DbDriver = "MySQL"
//var DbUser = "root"
//var DbPass = "19900718qzyQZY"
//var DbName = "barley"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Version       string `mapstructure:"VERSION"`
}

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
