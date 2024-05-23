package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	auth_pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jcchavezs/openfga-envoy/extauthz/internal/extractor"
	"github.com/jcchavezs/openfga-envoy/extauthz/internal/server/authz"
	"github.com/jcchavezs/openfga-envoy/extauthz/internal/server/config"
	"github.com/openfga/go-sdk/client"
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

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:               cfg.Server.APIURL,
		StoreId:              cfg.Server.StoreID,
		AuthorizationModelId: cfg.Server.AuthorizationModelID, // optional, recommended to be set for production
	})
	if err != nil {
		log.Fatalf("failed to initialize OpenFGA client: %v", err)
	}

	extractionSet := make([]extractor.ExtractorSet, 0, len(cfg.ExtractionSet))
	for _, es := range cfg.ExtractionSet {
		var (
			eSet extractor.ExtractorSet
			err  error
		)

		eSet.Name = es.Name

		eSet.User, err = extractor.MakeExtractor(es.User.Type, es.User.Config)
		if err != nil {
			log.Fatalf("failed to create user extractor: %v", err)
		}

		eSet.Object, err = extractor.MakeExtractor(es.Object.Type, es.Object.Config)
		if err != nil {
			log.Fatalf("failed to create object extractor: %v", err)
		}

		eSet.Relation, err = extractor.MakeExtractor(es.Relation.Type, es.Relation.Config)
		if err != nil {
			log.Fatalf("failed to create relation extractor: %v", err)
		}

		extractionSet = append(extractionSet, eSet)
	}

	filter := authz.NewExtAuthZFilter(fgaClient, extractionSet)

	server := createServer(filter)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting server on port %d\n", port)
	log.Fatal(server.Serve(listener))
}

func createServer(filter *authz.ExtAuthZFilter) *grpc.Server {
	grpcServer := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	auth_pb.RegisterAuthorizationServer(grpcServer, filter)

	return grpcServer
}
