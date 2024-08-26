package authz

import (
	"context"
	"errors"
	"fmt"
	"log"

	envoy "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	// Response for a successful authorization.
	allow = &envoy.CheckResponse{
		Status: &status.Status{
			Code:    int32(codes.OK),
			Message: "",
		},
	}

	deny = func(code codes.Code, message string) *envoy.CheckResponse {
		return &envoy.CheckResponse{
			Status: &status.Status{
				Code:    int32(code),
				Message: message,
			},
		}
	}
)

// ExtAuthZFilter is an implementation of the Envoy AuthZ filter.
type ExtAuthZFilter struct {
	client        *client.OpenFgaClient
	extractionKit []extractor.ExtractorKit
	modelID       string
}

var _ envoy.AuthorizationServer = (*ExtAuthZFilter)(nil)

// NewExtAuthZFilter creates a new ExtAuthZFilter
func NewExtAuthZFilter(c *client.OpenFgaClient, es []extractor.ExtractorKit) *ExtAuthZFilter {
	return &ExtAuthZFilter{client: c, extractionKit: es}
}

func (e ExtAuthZFilter) Register(server *grpc.Server) {
	envoy.RegisterAuthorizationServer(server, e)
}

func (e ExtAuthZFilter) Check(ctx context.Context, req *envoy.CheckRequest) (response *envoy.CheckResponse, err error) {
	res, err := e.check(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// TODO: replace with logging library
	log.Println(res)
	return res, nil
}

func (e ExtAuthZFilter) extract(ctx context.Context, req *envoy.CheckRequest) (*extractor.Check, error) {
	for _, es := range e.extractionKit {
		check, err := es.Extract(ctx, req)
		if err == nil {
			return check, nil
		}

		if errors.Is(err, extractor.ErrValueNotFound) {
			continue
		}

		return nil, err
	}

	return nil, nil
}

// Check implements the Check method of the Authorization interface.
func (e ExtAuthZFilter) check(ctx context.Context, req *envoy.CheckRequest) (response *envoy.CheckResponse, err error) {
	extracted, err := e.extract(ctx, req)
	if err != nil {
		return nil, err
	}

	if extracted == nil {
		return deny(codes.InvalidArgument, "No extraction set found"), nil
	}

	body := client.ClientCheckRequest{
		User:     extracted.User,
		Relation: extracted.Relation,
		Object:   extracted.Object,
		Context:  &extracted.Context,
	}

	options := client.ClientCheckOptions{
		AuthorizationModelId: openfga.PtrString(e.modelID),
	}

	data, err := e.client.Check(ctx).Body(body).Options(options).Execute()
	if err != nil {
		return deny(codes.Internal, fmt.Sprintf("Error checking permissions: %v", err)), nil
	}

	if data.GetAllowed() {
		return allow, nil
	}

	return deny(codes.PermissionDenied, fmt.Sprintf("Access denied: %s", data.GetResolution())), nil
}
