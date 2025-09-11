package databases

import (
	"context"

	"github.com/adrisongomez/thesis/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDatabase struct {
	cfg    *config.DatabaseConfig
	driver neo4j.DriverWithContext
}

func (db *Neo4jDatabase) Close(ctx context.Context) error {
	err := db.driver.Close(ctx)

	if err != nil {
		return err
	}
	db.driver = nil
	return nil
}

func (db *Neo4jDatabase) Connect(ctx context.Context) error {
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

func (db *Neo4jDatabase) GetDriver() neo4j.DriverWithContext {
	return db.driver
}
