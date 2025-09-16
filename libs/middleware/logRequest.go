package middleware

import (
	"fmt"

	"github.com/adrisongomez/thesis/libs/opentelemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func LogRequest(tp opentelemetry.TelemetryProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		tracer := tp.GetTracerProvider().Tracer(tp.GetServiceName())
		propagator := otel.GetTextMapPropagator()
		return func(c echo.Context) error {
			path_value := c.Request().URL.Path

			if path_value == "/api/healthcheck" {
				return next(c)
			}
			// Extract the context from the incoming request's headers.
			ctx := propagator.Extract(c.Request().Context(), propagation.HeaderCarrier(c.Request().Header))

			// Start a new span.
			spanName := fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path)
			ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
			defer span.End() // Ensure the span is closed when the function exits.

			// Set attributes on the span for better visibility.
			span.SetAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.url", c.Request().URL.String()),
				attribute.String("http.route", c.Path()),
			)

			// Inject the new context (with the span) into the request.
			c.SetRequest(c.Request().WithContext(ctx))

			// Call the next handler in the chain.
			err := next(c)

			// Set the span status based on the HTTP response code.
			status := c.Response().Status
			span.SetAttributes(attribute.Int("http.status_code", status))

			return err
		}
	}
}
