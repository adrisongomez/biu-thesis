package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	e := echo.New()
	e.Use(otelecho.Middleware("experiment-1"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{"*"}}))
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
