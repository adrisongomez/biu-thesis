package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"

	_ "google.golang.org/grpc/encoding/gzip" // Import to handle gzip compression

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	traceService "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	common "go.opentelemetry.io/proto/otlp/common/v1"
	resource "go.opentelemetry.io/proto/otlp/resource/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// SpanNode represents a single span in a format suitable for a Neo4j graph node.
// The "json" tags are included for easy marshaling if you want to log the struct as JSON.
type SpanNode struct {
	TraceID           string            `json:"traceId"`
	SpanID            string            `json:"spanId"`
	ParentSpanID      string            `json:"parentSpanId,omitempty"`
	Name              string            `json:"name"`
	Kind              string            `json:"kind"`
	ServiceName       string            `json:"serviceName"`
	StartTimeUnixNano uint64            `json:"startTimeUnixNano"`
	EndTimeUnixNano   uint64            `json:"endTimeUnixNano"`
	DurationNanos     uint64            `json:"durationNanos"`
	Attributes        map[string]string `json:"attributes"`
}

// ServiceNode represents a service that emits spans.
type ServiceNode struct {
	Name string `json:"name"`
}

// TraceNode represents a single trace, which is a collection of spans.
type TraceNode struct {
	TraceID string `json:"traceId"`
}

// Neo4jRepository handles all interactions with the Neo4j database.
type Neo4jRepository struct {
	driver neo4j.DriverWithContext
}

// NewNeo4jRepository creates a new repository and connects to the database.
func NewNeo4jRepository(ctx context.Context, uri, user, password string) (*Neo4jRepository, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, fmt.Errorf("could not create neo4j driver: %w", err)
	}
	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, fmt.Errorf("could not verify neo4j connectivity: %w", err)
	}
	log.Println("Successfully connected to Neo4j.")
	return &Neo4jRepository{driver: driver}, nil
}

// Close closes the database connection.
func (r *Neo4jRepository) Close(ctx context.Context) {
	r.driver.Close(ctx)
}

// SaveBatch saves spans, services, and traces to Neo4j in a single transaction.
func (r *Neo4jRepository) SaveBatch(ctx context.Context, spans []SpanNode, services map[string]ServiceNode, traces map[string]TraceNode) error {
	if len(spans) == 0 {
		return nil
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Convert structs to maps for the driver
	spanMaps := make([]map[string]interface{}, len(spans))
	for i, s := range spans {
		spanMaps[i] = map[string]interface{}{
			"traceId":           s.TraceID,
			"spanId":            s.SpanID,
			"parentSpanId":      s.ParentSpanID,
			"name":              s.Name,
			"kind":              s.Kind,
			"serviceName":       s.ServiceName,
			"startTimeUnixNano": s.StartTimeUnixNano,
			"endTimeUnixNano":   s.EndTimeUnixNano,
			"durationNanos":     s.DurationNanos,
			"attributes":        s.Attributes,
		}
	}
	serviceMaps := make([]map[string]interface{}, 0, len(services))
	for _, s := range services {
		serviceMaps = append(serviceMaps, map[string]interface{}{"name": s.Name})
	}
	traceMaps := make([]map[string]interface{}, 0, len(traces))
	for _, t := range traces {
		traceMaps = append(traceMaps, map[string]interface{}{"traceId": t.TraceID})
	}

	// Use ExecuteWrite for a managed transaction.
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// Using multiple queries in a single transaction is often clearer and safer
		// than one massive, complex query.

		// 1. Create Service nodes
		if _, err := tx.Run(ctx, `UNWIND $services AS service_props MERGE (s:Service {name: service_props.name})`, map[string]interface{}{"services": serviceMaps}); err != nil {
			return nil, err
		}
		// 2. Create Trace nodes
		if _, err := tx.Run(ctx, `UNWIND $traces AS trace_props MERGE (t:Trace {traceId: trace_props.traceId})`, map[string]interface{}{"traces": traceMaps}); err != nil {
			return nil, err
		}
		// 3. Create Span nodes
		if _, err := tx.Run(ctx, `UNWIND $spans AS span_props MERGE (s:Span {spanId: span_props.spanId}) SET s = span_props`, map[string]interface{}{"spans": spanMaps}); err != nil {
			return nil, err
		}
		// 4. Create relationships
		relQueries := `
		// Link Services to Spans
		UNWIND $spans AS span_props
		MATCH (service:Service {name: span_props.serviceName})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (service)-[:EMITTED]->(span);

		// Link Traces to Spans
		UNWIND $spans AS span_props
		MATCH (trace:Trace {traceId: span_props.traceId})
		MATCH (span:Span {spanId: span_props.spanId})
		MERGE (trace)-[:CONTAINS]->(span);

		// Link parent Spans to child Spans
		UNWIND $spans AS span_props
		WHERE span_props.parentSpanId IS NOT NULL AND span_props.parentSpanId <> ""
		MATCH (parent:Span {spanId: span_props.parentSpanId})
		MATCH (child:Span {spanId: span_props.spanId})
		MERGE (parent)-[:PARENT_OF]->(child);
		`
		// Using EagerResult() can help ensure one part of a script finishes before the next
		if _, err := tx.Run(ctx, relQueries, map[string]interface{}{"spans": spanMaps}); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute batch write transaction: %w", err)
	}

	log.Printf("Successfully wrote %d spans, %d services, and %d traces to Neo4j.", len(spans), len(services), len(traces))
	return nil
}

// traceServiceServer is our implementation of the OTLP TraceService gRPC service.
type traceServiceServer struct {
	traceService.UnimplementedTraceServiceServer
	repo *Neo4jRepository
}

// Export is the method that receives trace data from the OpenTelemetry Collector.
func (s *traceServiceServer) Export(ctx context.Context, req *traceService.ExportTraceServiceRequest) (*traceService.ExportTraceServiceResponse, error) {
	spansToSave := []SpanNode{}
	servicesToSave := make(map[string]ServiceNode)
	tracesToSave := make(map[string]TraceNode)

	resourceSpans := req.GetResourceSpans()
	for _, rs := range resourceSpans {
		serviceName := getServiceName(rs.GetResource())
		if _, ok := servicesToSave[serviceName]; !ok {
			servicesToSave[serviceName] = ServiceNode{Name: serviceName}
		}

		scopeSpans := rs.GetScopeSpans()
		for _, ss := range scopeSpans {
			for _, span := range ss.GetSpans() {
				traceID := hex.EncodeToString(span.GetTraceId())
				if _, ok := tracesToSave[traceID]; !ok {
					tracesToSave[traceID] = TraceNode{TraceID: traceID}
				}

				spanNode := SpanNode{
					TraceID:           traceID,
					SpanID:            hex.EncodeToString(span.GetSpanId()),
					ParentSpanID:      hex.EncodeToString(span.GetParentSpanId()),
					Name:              span.GetName(),
					Kind:              span.GetKind().String(),
					ServiceName:       serviceName,
					StartTimeUnixNano: span.GetStartTimeUnixNano(),
					EndTimeUnixNano:   span.GetEndTimeUnixNano(),
					DurationNanos:     span.GetEndTimeUnixNano() - span.GetStartTimeUnixNano(),
					Attributes:        attributesToMap(span.GetAttributes()),
				}
				spansToSave = append(spansToSave, spanNode)
			}
		}
	}

	if len(spansToSave) > 0 {
		if err := s.repo.SaveBatch(ctx, spansToSave, servicesToSave, tracesToSave); err != nil {
			log.Printf("Error saving batch to Neo4j: %v", err)
		}
	}

	return &traceService.ExportTraceServiceResponse{}, nil
}

// getServiceName extracts the service.name attribute from the resource.
func getServiceName(res *resource.Resource) string {
	for _, attr := range res.GetAttributes() {
		if attr.GetKey() == "service.name" {
			return attr.GetValue().GetStringValue()
		}
	}
	return "unknown-service"
}

// attributesToMap converts a slice of OTLP KeyValue attributes to a simple map.
func attributesToMap(attrs []*common.KeyValue) map[string]string {
	m := make(map[string]string)
	for _, attr := range attrs {
		m[attr.GetKey()] = attr.GetValue().GetStringValue()
	}
	return m
}

func main() {
	ctx := context.Background()

	// --- Neo4j Connection Setup ---
	neo4jURI := os.Getenv("NEO4J_URI")
	if neo4jURI == "" {
		neo4jURI = "neo4j://localhost:7687" // Default value
	}
	neo4jUser := os.Getenv("NEO4J_USER")
	if neo4jUser == "" {
		neo4jUser = "neo4j" // Default value
	}
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	if neo4jPassword == "" {
		log.Fatal("NEO4J_PASSWORD environment variable not set")
	}

	repo, err := NewNeo4jRepository(ctx, neo4jURI, neo4jUser, neo4jPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer repo.Close(ctx)

	// --- gRPC Server Setup ---
	port := 4317
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register our trace service implementation, passing the repository.
	traceService.RegisterTraceServiceServer(s, &traceServiceServer{repo: repo})

	reflection.Register(s)

	log.Printf("OTLP trace receiver listening on gRPC port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
