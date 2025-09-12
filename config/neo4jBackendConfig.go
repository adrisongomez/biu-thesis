package config

import "github.com/spf13/viper"

type Neo4jBackendConfig struct {
	ServicePort string `mapstrucutre:"SERVICE_PORT" default:"4317"`
	DbUserName  string `mapstructure:"DB_USERNAME"`
	DbPassword  string `mapstructure:"DB_PASSWORD"`
	DbUrl       string `mapstructure:"DB_URI"`
}

func NewDatabaseConfig(pathToConfig string) (*Neo4jBackendConfig, error) {
	viper.SetConfigFile(pathToConfig)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &Neo4jBackendConfig{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
