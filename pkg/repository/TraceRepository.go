package repository

import (
	"context"

	"github.com/adrisongomez/thesis/pkg/models"
)

type TraceRepository interface {
	SaveBatch(ctx context.Context, spans []models.SpanNode, services map[string]models.ServiceNode, traces map[string]models.TraceNode) error
}
