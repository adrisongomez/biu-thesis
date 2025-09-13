package databases

import (
	"context"

	"github.com/adrisongomez/thesis/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

type Neo4jConnector struct {
	logger *zap.SugaredLogger
	cfg    *config.Neo4jBackendConfig
	driver neo4j.DriverWithContext
}

func (db *Neo4jConnector) Close(ctx context.Context) error {
	err := db.driver.Close(ctx)

	if err != nil {
		return err
	}
	db.driver = nil
	return nil
}

func (db *Neo4jConnector) Connect(ctx context.Context) error {
	logger := db.logger
	logger.Infow("Neo4jConnector#Connect got called", "uri", db.cfg.DbUrl, "username", db.cfg.DbUserName, "pw", db.cfg.DbPassword)
	driver, err := neo4j.NewDriverWithContext(db.cfg.DbUrl, neo4j.BasicAuth(db.cfg.DbUserName, db.cfg.DbPassword, ""))
	if err != nil {
		logger.Warn("Error trying to connect with database")
		return err
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		logger.Warn("Error verifting connectivity")
		return err
	}
	db.driver = driver
	return nil
}

func (db *Neo4jConnector) GetDriver() neo4j.DriverWithContext {
	db.logger.Info("Neo4jConnector#GetDriver got called")
	return db.driver
}

func NewNeo4jConnector(cfg *config.Neo4jBackendConfig) *Neo4jConnector {
	logger := zap.L().Sugar()
	logger.Infow("NewNeo4jConnector Got Called with", "url", cfg.DbUrl, "servicePort", cfg.ServicePort)
	return &Neo4jConnector{
		logger: logger,
		cfg:    cfg,
		driver: nil,
	}
}
