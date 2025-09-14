package main

import (
	"context"
	"log"

	"github.com/adrisongomez/thesis/config"
	"github.com/adrisongomez/thesis/libs/databases"
	"github.com/adrisongomez/thesis/libs/loggers"
	"github.com/adrisongomez/thesis/pkg/repository"
	"github.com/adrisongomez/thesis/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()
	logger, err := loggers.NewLogger()

	if err != nil {
		log.Fatalf("Error trying to initialize the logger %v", err)
		return
	}

	cfg, err := config.NewDatabaseConfig(".env")

	if err != nil {
		logger.Errorw("Error trying to load the configuration", "error", err)
		return
	}

	conn := databases.NewNeo4jConnector(cfg)
	conn.Connect(ctx)
	defer conn.Close(ctx)
	repo := repository.NewNeo4jTraceQueryRepository(conn)
	service := services.NewTraceQueryService(repo)

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{"*"}}))
	e.GET("/api/traces/:traceId/spans", service.GetTraceHandler)
	e.GET("/api/traces", service.GetTracesList)

	if err := e.Start(":" + cfg.ServicePort); err != nil {
		logger.Errorw("failed to start API server", "error", err)
	}
}
