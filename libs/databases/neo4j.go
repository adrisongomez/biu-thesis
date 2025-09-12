package databases

import (
	"context"

	"github.com/adrisongomez/thesis/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jConnector struct {
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
	driver, err := neo4j.NewDriverWithContext(db.cfg.DbUrl, neo4j.BasicAuth(db.cfg.DbUserName, db.cfg.DbPassword, ""))
	if err != nil {
		return err
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	}
	db.driver = driver
	return nil
}

func (db *Neo4jConnector) GetDriver() neo4j.DriverWithContext {
	return db.driver
}

func NewNeo4jConnector(cfg *config.Neo4jBackendConfig) *Neo4jConnector {
	return &Neo4jConnector{
		cfg:    cfg,
		driver: nil,
	}
}
