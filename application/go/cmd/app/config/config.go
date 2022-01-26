package config

import "github.com/spf13/viper"

var (
	Version   = "UNKNOWN"
	BuildDate = "UNKNOWN"
)

type Config struct {
	ApiBasePath               string
	PsqlConnectionStringRead  string
	PsqlConnectionStringWrite string
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

	viper.SetDefault("ApiBasePath", "/api/")
	viper.SetDefault("PsqlConnectionStringRead", "postgresql://dbname:dbpw@localhost:5432/dbname")
	viper.SetDefault("PsqlConnectionStringWrite", "postgresql://dbname:dbpw@localhost:5432/dbname")

	return viper.ReadInConfig()
}

func parseViperConfig() (*Config, error) {
	return &Config{
		ApiBasePath:               viper.GetString("ApiBasePath"),
		PsqlConnectionStringRead:  viper.GetString("PsqlConnectionStringRead"),
		PsqlConnectionStringWrite: viper.GetString("PsqlConnectionStringWrite"),
	}, nil
}
