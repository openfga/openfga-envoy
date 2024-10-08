package main

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	auth_pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"github.com/openfga/openfga-envoy/extauthz/internal/server/authz"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context, e extractor.ExtractorKit) (auth_pb.AuthorizationClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: "https://api.fga.example",
	})
	if err != nil {
		panic(err)
	}

	filter := authz.NewExtAuthZFilter(authz.Config{
		Enforce:        true,
		ExtractionKits: []extractor.ExtractorKit{e},
	}, fgaClient, logger.NewNoopLogger())

	baseServer := grpc.NewServer()
	auth_pb.RegisterAuthorizationServer(baseServer, filter)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	return auth_pb.NewAuthorizationClient(conn), closer
}

func TestNoUserExtractedFails(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("no user")

	e := extractor.ExtractorKit{
		Name: "extauthz",
		User: func(ctx context.Context, value *auth_pb.CheckRequest) (extractor.Extraction, bool, error) {
			return extractor.Extraction{}, false, expectedErr
		},
	}

	client, closer := server(ctx, e)
	defer closer()

	_, sErr := client.Check(ctx, &auth_pb.CheckRequest{})
	if sErr == nil {
		t.Fatal("expected error")
	}

	require.ErrorContains(t, sErr, expectedErr.Error())
}
