package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/adrisongomez/thesis/libs/databases"
	"github.com/adrisongomez/thesis/pkg/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

type Neo4jTraceQueryRepository struct {
	lg   *zap.SugaredLogger
	conn *databases.Neo4jConnector
}

func (tq *Neo4jTraceQueryRepository) GetTraceByID(ctx context.Context, traceID string) ([]models.SpanNode, error) {
	tq.lg.Infow("Neo4jTraceQueryRepository#GetTraceByID", "traceID", traceID)
	session := tq.conn.GetDriver().NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
		MATCH (s:Span {traceId: $traceId})
		RETURN s
		ORDER BY s.startTimeUnixNano ASC`
		records, err := tx.Run(ctx, query, map[string]interface{}{"traceId": traceID})
		if err != nil {
			return nil, err
		}

		spans := []models.SpanNode{}
		// Define core fields to distinguish them from custom attributes
		coreFields := map[string]struct{}{
			"traceId": {}, "spanId": {}, "parentSpanId": {},
			"name": {}, "kind": {}, "serviceName": {},
			"startTimeUnixNano": {}, "endTimeUnixNano": {}, "durationNanos": {},
		}

		for records.Next(ctx) {
			record := records.Record()
			node, ok := record.Values[0].(neo4j.Node)
			if !ok {
				return nil, fmt.Errorf("could not cast to neo4j.Node")
			}
			props := node.Props
			attributes := make(map[string]string)

			// Reconstruct attributes map by iterating over all node properties
			for key, value := range props {
				if _, isCore := coreFields[key]; !isCore {
					if strVal, isStr := value.(string); isStr {
						// This is a best-effort "un-sanitize". For perfect mapping,
						// you might store original attribute keys in a separate property.
						originalKey := strings.ReplaceAll(key, "_", ".")
						attributes[originalKey] = strVal
					}
				}
			}

			parentSpanID := ""
			if p, ok := props["parentSpanId"].(string); ok {
				parentSpanID = p
			}

			span := models.SpanNode{
				TraceID:           props["traceId"].(string),
				SpanID:            props["spanId"].(string),
				ParentSpanID:      parentSpanID,
				Name:              props["name"].(string),
				Kind:              props["kind"].(string),
				ServiceName:       props["serviceName"].(string),
				StartTimeUnixNano: uint64(props["startTimeUnixNano"].(int64)),
				EndTimeUnixNano:   uint64(props["endTimeUnixNano"].(int64)),
				DurationNanos:     uint64(props["durationNanos"].(int64)),
				Attributes:        attributes,
			}
			spans = append(spans, span)
		}

		return spans, records.Err()
	})
	if err != nil {
		return nil, err
	}
	response, ok := result.([]models.SpanNode)
	if !ok {
		tq.lg.Errorw("Error casting spandNode", "result", result)
		return nil, fmt.Errorf("failed to cast result to []models.SpanNode")
	}
	tq.lg.Infow("Neo4jTraceQueryRepository#GetTraceByID returns", "spans", len(response))
	return response, nil
}

func (r *Neo4jTraceQueryRepository) GetTraces(ctx context.Context) ([]models.TraceSummary, error) {
	session := r.conn.GetDriver().NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			// Find all traces
			MATCH (t:Trace)
			// For each trace, find all its spans
			MATCH (t)-[:CONTAINS]->(s:Span)
			// For each trace, find its root span (a span with no incoming PARENT_OF relationship)
			MATCH (root:Span {traceId: t.traceId})
			WHERE NOT (root)<-[:PARENT_OF]-(:Span)
			// Get the service that emitted the root span
			MATCH (service:Service)-[:EMITTED]->(root)
			// Aggregate the data per trace
			WITH t,
				 root,
				 service,
				 min(s.startTimeUnixNano) as startTime,
				 max(s.endTimeUnixNano) as endTime,
				 count(s) as spanCount
			// Order by the most recent traces first
			ORDER BY startTime DESC
			// Limit the result set
			LIMIT 100
			// Return the summary data
			RETURN t.traceId as traceId,
				   root.name as rootSpanName,
				   service.name as serviceName,
				   startTime,
				   (endTime - startTime) as duration,
				   spanCount`
		records, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		summaries := []models.TraceSummary{}
		for records.Next(ctx) {
			record := records.Record()
			summary := models.TraceSummary{
				TraceID:      record.Values[0].(string),
				RootSpanName: record.Values[1].(string),
				ServiceName:  record.Values[2].(string),
				StartTime:    uint64(record.Values[3].(int64)),
				Duration:     uint64(record.Values[4].(int64)),
				SpanCount:    record.Values[5].(int64),
			}
			summaries = append(summaries, summary)
		}
		return summaries, records.Err()
	})
	if err != nil {
		return nil, err
	}
	return result.([]models.TraceSummary), nil
}

func NewNeo4jTraceQueryRepository(conn *databases.Neo4jConnector) TraceQueryRepository {
	lg := zap.L().Sugar()
	return &Neo4jTraceQueryRepository{
		conn: conn,
		lg:   lg,
	}
}
