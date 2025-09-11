package main

import (
	"context"
	"log"
	"net/http"

	m "github.com/adrisongomez/thesis/libs/middleware"
	"github.com/adrisongomez/thesis/config"
	"github.com/adrisongomez/thesis/libs/opentelemetry"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig("./.env")
	if err != nil {
		log.Fatalf("error reading configuration %d", err)
		return
	}
	t, err := opentelemetry.NewTelemetry(ctx, cfg)
	if err != nil {
		log.Fatalf("Error starting the opentelemetry config %d", err)
		return
	}
	defer t.Shutdown(ctx)
	e := echo.New()
	e.Use(otelecho.Middleware(t.GetServiceName(), otelecho.WithTracerProvider(t.GetTracerProvider())))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{"*"}}))
	e.Use(m.LogRequest(t))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.POST("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.PUT("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.DELETE("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3003"))
}
