package config

import (
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/spf13/viper"
)

var (
	Version   = "UNKNOWN"
	BuildDate = "UNKNOWN"
)

type Config struct {
	SessionKey    string
	UsersLoginMap connection.UsersLoginMap
}

func GenerateConfig() (*Config, error) {
	e := setupDefaultViperConfig()
	if e != nil {
		return nil, e
	}
	return parseViperConfig()
}

func setupDefaultViperConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf/")
	viper.AddConfigPath(".")

	viper.SetDefault("SessionKey", "your-secret-key")

	usersLoginMap := make(connection.UsersLoginMap)
	usersLoginMap["supplier1"] = "pswSupplier1"
	viper.SetDefault("UsersLoginMap", usersLoginMap)

	return viper.ReadInConfig()
}

func parseViperConfig() (*Config, error) {
	usersLoginMap := make(connection.UsersLoginMap)
	err := viper.UnmarshalKey("UsersLoginMap", &usersLoginMap)
	if err != nil {
		return nil, err
	}

	return &Config{
		SessionKey:    viper.GetString("SessionKey"),
		UsersLoginMap: usersLoginMap,
	}, nil
}
