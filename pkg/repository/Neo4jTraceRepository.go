package repository

import (
	"context"
	"fmt"

	"github.com/adrisongomez/thesis/libs/databases"
	"github.com/adrisongomez/thesis/pkg/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

type Neo4jTraceRepository struct {
	conn   *databases.Neo4jConnector
	logger *zap.SugaredLogger
}

func (r *Neo4jTraceRepository) SaveBatch(ctx context.Context, spans []models.SpanNode, services map[string]models.ServiceNode, traces map[string]models.TraceNode) error {
	r.logger.Infow("NeoTraceRespository#SaveBatch got called with", "spans", spans, "services", services, "traces", traces)
	driver := r.conn.GetDriver()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})

	rawSpans := make([]map[string]any, len(spans))
	for i, s := range spans {
		rawSpans[i] = s.Raw()
	}

	rawServices := make([]map[string]any, 0, len(services))
	for _, s := range services {
		rawServices = append(rawServices, s.Raw())
	}

	rawTraces := make([]map[string]any, 0, len(traces))
	for _, t := range traces {
		rawTraces = append(rawTraces, t.Raw())
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if _, err := tx.Run(ctx, CREATE_SERVICES, map[string]any{"services": rawServices}); err != nil {
			r.logger.Error("CreateService got an error", err)
			return nil, err
		}
		if _, err := tx.Run(ctx, CREATE_TRACES, map[string]any{"traces": rawTraces}); err != nil {
			r.logger.Error("CreateTrace got an error", err)
			return nil, err
		}
		if _, err := tx.Run(ctx, CREATE_SPANS, map[string]any{"spans": rawSpans}); err != nil {
			r.logger.Error("CreateSpans got an error", err)
			return nil, err
		}
		if _, err := tx.Run(ctx, CREATE_RELATIONSHIPS, map[string]any{"spans": rawSpans}); err != nil {
			r.logger.Error("CreateRelationships got an error", err)
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		r.logger.Error("Neo4j#ExecuteWrite got an error", err)
		return fmt.Errorf("failed to execute batch write transaction: %w", err)
	}
	r.logger.Info("Successfully wrote %d spans, %d services, and %d traces to Neo4j.", len(spans), len(services), len(traces))
	return nil
}

func NewNeo4jTraceRepository(conn *databases.Neo4jConnector) *Neo4jTraceRepository {
	logger := zap.L().Sugar()
	logger.Info("NewNeo4jTraceRepository got called")
	return &Neo4jTraceRepository{conn: conn, logger: logger}
}
