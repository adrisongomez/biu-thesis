package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Neo4jBackendConfig struct {
	ServicePort string `mapstrucutre:"SERVICE_PORT" default:"4317"`
	DbUserName  string `mapstructure:"DB_USERNAME"`
	DbPassword  string `mapstructure:"DB_PASSWORD"`
	DbUrl       string `mapstructure:"DB_URI"`
}

func NewDatabaseConfig(pathToConfig string) (*Neo4jBackendConfig, error) {
	logger := zap.L().Sugar()
	logger.Infof("New Database config got called with %s", pathToConfig)
	viper.SetConfigFile(pathToConfig)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Warnf("Some error got caught at loading the database config %s", pathToConfig)
		return nil, err
	}
	cfg := &Neo4jBackendConfig{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		logger.Warnf("Some error got caught unmarshing the database config %s", pathToConfig)
		return nil, err
	}
	return cfg, nil
}
