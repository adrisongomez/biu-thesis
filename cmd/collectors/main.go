package main

import (
	"fmt"
	"log"
	"net"

	"github.com/adrisongomez/thesis/config"
	"github.com/adrisongomez/thesis/libs/databases"
	"github.com/adrisongomez/thesis/libs/loggers"
	"github.com/adrisongomez/thesis/pkg/repository"
	"github.com/adrisongomez/thesis/pkg/services"
	traceService "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	logger, err := loggers.NewLogger()
	if err != nil {
		log.Fatalf("Error trying to initilizing the logger: %v", err)
		return
	}
	defer logger.Sync()
	cfg, err := config.NewDatabaseConfig(".env")
	if err != nil {
		logger.Errorf("Error trying to load the environment variable from config file %v", err)
		return
	}

	conn := databases.NewNeo4jConnector(cfg)
	err = conn.Connect(ctx)
	if err != nil {
		logger.Errorf("Error while trying to connect with Neo4j Database %v", err)
		return
	}
	defer conn.Close(ctx)

	s := grpc.NewServer()
	repo := repository.NewNeo4jTraceRepository(conn)
	traceService.RegisterTraceServiceServer(s, services.NewTraceServiceServer(repo))
	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServicePort))

	if err != nil {
		logger.Errorf("Error trying to open listener for port %s - %v", cfg.ServicePort, err)
		return
	}

	logger.Infof("OTLP trace receiver listening on gRPC port %s", cfg.ServicePort)
	if err := s.Serve(lis); err != nil {
		logger.Errorf("failed to serve: %v", err)
	}
}
