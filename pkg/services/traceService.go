package services

import (
	"context"
	"log"

	"github.com/adrisongomez/thesis/pkg/models"
	"github.com/adrisongomez/thesis/pkg/repository"
	traceService "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	resource "go.opentelemetry.io/proto/otlp/resource/v1"
	"go.uber.org/zap"
	_ "google.golang.org/grpc/encoding/gzip"
)

type TraceServiceServer struct {
	lg   *zap.SugaredLogger
	repo repository.TraceRepository
	traceService.UnimplementedTraceServiceServer
}

func (s *TraceServiceServer) Export(ctx context.Context, req *traceService.ExportTraceServiceRequest) (*traceService.ExportTraceServiceResponse, error) {
	s.lg.Infow("TraceServiceServer#Export got called with", "req", req.String())
	spansToSave := []models.SpanNode{}
	servicesToSave := make(map[string]models.ServiceNode)
	tracesToSave := make(map[string]models.TraceNode)
	s.lg.Info("TraceServiceServer#Export extracting data from req")
	for _, rs := range req.GetResourceSpans() {
		serviceName := getServiceName(rs.GetResource())
		if _, ok := servicesToSave[serviceName]; !ok {
			servicesToSave[serviceName] = models.NewServiceNode(serviceName)
		}

		for _, ss := range rs.GetScopeSpans() {
			for _, span := range ss.GetSpans() {
				spanNode := models.NewSpanNodeFromV1Span(span, serviceName)
				if _, ok := tracesToSave[spanNode.TraceID]; !ok {
					tracesToSave[spanNode.TraceID] = models.NewTraceNode(spanNode.TraceID)
				}
				spansToSave = append(spansToSave, spanNode)
			}
		}
	}

	s.lg.Info("TraceServiceServer#Export calling repository to save the data")
	if len(spansToSave) > 0 {
		if err := s.repo.SaveBatch(ctx, spansToSave, servicesToSave, tracesToSave); err != nil {
			log.Printf("Error saving batch to Neo4j: %v", err)
		}
	}
	return &traceService.ExportTraceServiceResponse{}, nil
}

func NewTraceServiceServer(repo repository.TraceRepository) *TraceServiceServer {
	logger := zap.L().Sugar()
	logger.Info("NewTraceServiceServer got called with")
	return &TraceServiceServer{repo: repo, lg: logger}
}

func getServiceName(res *resource.Resource) string {
	logger := zap.L().Sugar()
	logger.Info("getServiceName got called")
	for _, attr := range res.GetAttributes() {
		if attr.GetKey() == "service.name" {
			return attr.GetValue().GetStringValue()
		}
	}
	return "unknown-service"
}
