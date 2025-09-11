package opentelemetry

import (
	"fmt"

	"github.com/adrisongomez/thesis/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/context"
)

type TelemetryProvider interface {
	GetServiceName() string
	GetTracerProvider() *trace.TracerProvider
	GetResource() *resource.Resource
	Shutdown(ctx context.Context)
}

type Telemetry struct {
	lp  *log.LoggerProvider
	tp  *trace.TracerProvider
	mp  *metric.MeterProvider
	cfg *config.TelemetryConfig
	res *resource.Resource
}

func (t *Telemetry) GetServiceName() string {
	return t.cfg.ServiceName
}

func (t *Telemetry) Shutdown(ctx context.Context) {
	t.lp.Shutdown(ctx)
	t.mp.Shutdown(ctx)
	t.tp.Shutdown(ctx)
}

func (t *Telemetry) GetTracerProvider() *trace.TracerProvider {
	return t.tp
}

func (t *Telemetry) GetResource() *resource.Resource {
	return t.res
}

func NewTelemetry(ctx context.Context, cfg *config.TelemetryConfig) (TelemetryProvider, error) {
	res := newResource(cfg)
	lp, err := newLoggerProvider(ctx, cfg, res)

	if err != nil {
		return nil, fmt.Errorf("failed to create logger %w", err)
	}

	mp, err := newMeterProvider(ctx, cfg, res)

	if err != nil {
		return nil, fmt.Errorf("failed to create metric %w", err)
	}

	// meter := mp.Meter(cfg.ServicName)

	tp, err := newTracerProvider(ctx, cfg, res)

	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return &Telemetry{
		lp:  lp,
		mp:  mp,
		tp:  tp,
		cfg: cfg,
		res: res,
	}, nil
}
