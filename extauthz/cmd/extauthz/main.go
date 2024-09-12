package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	auth_pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"github.com/openfga/openfga-envoy/extauthz/internal/server/authz"
	"github.com/openfga/openfga-envoy/extauthz/internal/server/config"
	"github.com/openfga/openfga/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const port = 9002

func main() {
	var (
		configPath string
	)

	flag.StringVar(&configPath, "config", "./config.yaml", "path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logger.NewLogger(parseLogConfig(cfg.Log)...)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	var filter auth_pb.AuthorizationServer = authz.NoopFilter{}
	if cfg.Mode != config.AuthModeDisabled {
		fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
			ApiUrl:               cfg.Server.APIURL,
			StoreId:              cfg.Server.StoreID,
			AuthorizationModelId: cfg.Server.AuthorizationModelID, // optional, recommended to be set for production
		})
		if err != nil {
			logger.Fatal("failed to initialize OpenFGA client", zap.Error(err))
		}

		extractionSet := make([]extractor.ExtractorKit, 0, len(cfg.ExtractionSet))
		for _, es := range cfg.ExtractionSet {
			var (
				eSet extractor.ExtractorKit
				err  error
			)

			eSet.Name = es.Name

			eSet.User, err = extractor.MakeExtractor(es.User.Type, es.User.Config)
			if err != nil {
				logger.Fatal("failed to create user extractor", zap.Error(err))
			}

			eSet.Object, err = extractor.MakeExtractor(es.Object.Type, es.Object.Config)
			if err != nil {
				logger.Fatal("failed to create object extractor", zap.Error(err))
			}

			eSet.Relation, err = extractor.MakeExtractor(es.Relation.Type, es.Relation.Config)
			if err != nil {
				logger.Fatal("failed to create relation extractor", zap.Error(err))
			}

			extractionSet = append(extractionSet, eSet)
		}

		filter = authz.NewExtAuthZFilter(
			authz.Config{
				Enforce:        cfg.Mode == config.AuthModeEnforce,
				ExtractionKits: extractionSet,
			},
			fgaClient,
			logger,
		)
	}

	server := createServer(filter)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed start listener: %v", err)
	}

	logger.Info("Starting server", zap.Int("port", port), zap.String("mode", cfg.Mode.String()))
	log.Fatal(server.Serve(listener))
}

func createServer(filter auth_pb.AuthorizationServer) *grpc.Server {
	grpcServer := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	auth_pb.RegisterAuthorizationServer(grpcServer, filter)

	return grpcServer
}

func parseLogConfig(cfg config.Log) []logger.OptionLogger {
	return []logger.OptionLogger{
		logger.WithLevel(cfg.Level),
		logger.WithFormat(cfg.Format),
		logger.WithTimestampFormat(cfg.TimestampFormat),
	}
}
