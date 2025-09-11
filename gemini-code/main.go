package main

import (
	"context"
	"fmt"
	"log"
	"net"

	_ "google.golang.org/grpc/encoding/gzip" // Import to handle gzip compression

	traceService "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	common "go.opentelemetry.io/proto/otlp/common/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// traceServiceServer is our implementation of the OTLP TraceService gRPC service.
// It must embed UnimplementedTraceServiceServer to satisfy the interface.
type traceServiceServer struct {
	traceService.UnimplementedTraceServiceServer
}

// Export is the method that receives trace data from the OpenTelemetry Collector.
func (s *traceServiceServer) Export(ctx context.Context, req *traceService.ExportTraceServiceRequest) (*traceService.ExportTraceServiceResponse, error) {
	log.Println("Received traces!")

	// The request contains a collection of ResourceSpans.
	resourceSpans := req.GetResourceSpans()
	log.Printf("Number of resource spans: %d", len(resourceSpans))

	// Iterate over each ResourceSpan
	for _, rs := range resourceSpans {
		// A ResourceSpan describes the entity that produced the spans (e.g., a service).
		resource := rs.GetResource()
		if resource != nil {
			log.Printf("  Resource: %s", formatAttributes(resource.GetAttributes()))
		}

		// Each ResourceSpan contains ScopeSpans (formerly InstrumentationLibrarySpans).
		scopeSpans := rs.GetScopeSpans()
		for _, ss := range scopeSpans {
			// A ScopeSpan is a collection of spans from a single instrumentation scope (e.g., a library).
			scope := ss.GetScope()
			if scope != nil {
				log.Printf("    Scope: %s", scope.GetName())
			}

			// Each ScopeSpan contains the actual spans.
			spans := ss.GetSpans()
			log.Printf("    Number of spans in this scope: %d", len(spans))
			for _, span := range spans {
				log.Printf("      - Span: %s (TraceID: %x, SpanID: %x, ParentSpanID: %x)",
					span.GetName(),
					span.GetTraceId(),
					span.GetSpanId(),
					span.GetParentSpanId(),
				)
				// Log span attributes
				if len(span.GetAttributes()) > 0 {
					log.Printf("        Attributes: %s", formatAttributes(span.GetAttributes()))
				}
			}
		}
	}

	// For this simple service, we always return a successful response.
	// In a real-world scenario, you might return partial success or an error.
	return &traceService.ExportTraceServiceResponse{}, nil
}

// formatAttributes is a helper function to convert OTLP KeyValue attributes to a readable string.
func formatAttributes(attrs []*common.KeyValue) string {
	var result string
	for i, attr := range attrs {
		result += fmt.Sprintf("%s=%s", attr.GetKey(), attr.GetValue().GetStringValue())
		if i < len(attrs)-1 {
			result += ", "
		}
	}
	return result
}

func main() {
	// OTLP standard gRPC port is 4317
	port := 4318
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register our trace service implementation.
	traceService.RegisterTraceServiceServer(s, &traceServiceServer{})

	// Register gRPC reflection service for tools like grpcurl to work.
	reflection.Register(s)

	log.Printf("OTLP trace receiver listening on gRPC port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
