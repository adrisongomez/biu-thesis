package repository

import (
	"context"

	"github.com/adrisongomez/thesis/pkg/models"
)

type (
	TraceRepository interface {
		SaveBatch(ctx context.Context, spans []models.SpanNode, services map[string]models.ServiceNode, traces map[string]models.TraceNode) error
	}
	TraceQueryRepository interface {
		GetTraceByID(ctx context.Context, traceID string) ([]models.SpanNode, error)
		GetTraces(ctx context.Context) ([]models.TraceSummary, error)
	}
)
