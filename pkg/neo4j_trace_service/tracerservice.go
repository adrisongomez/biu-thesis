package neo4jtraceservice

import (
	"context"
	"log"

	"github.com/adrisongomez/thesis/libs/databases"
	traceService "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type TraceServiceServer struct {
	db *databases.Neo4jDatabase
	traceService.UnimplementedTraceServiceServer
}

func (s *TraceServiceServer) Export(ctx context.Context, req *traceService.ExportTraceServiceRequest) (*traceService.ExportTraceServiceResponse, error) {
	err := s.db.Connect(ctx)
	if err != nil {
		return nil, err
	}
	spans := req.GetResourceSpans()
	for _, span := range spans {
		log.Println(span.String())
	}
	return nil, nil
}
