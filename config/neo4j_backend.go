package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	DbUserName string `mapstructure:"DB_USERNAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbUrl      string `mapstructure:"DB_URI"`
}

func NewDatabaseConfig(pathToConfig string) (*DatabaseConfig, error) {
	viper.SetConfigFile(pathToConfig)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &DatabaseConfig{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
