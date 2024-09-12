package authz

import (
	"context"
	"errors"
	"fmt"

	envoy "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"github.com/openfga/openfga/pkg/logger"
	"go.uber.org/zap"
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
	enforce        bool
	client         *client.OpenFgaClient
	extractionKits []extractor.ExtractorKit
	modelID        string
	logger         logger.Logger
}

var _ envoy.AuthorizationServer = (*ExtAuthZFilter)(nil)

type Config struct {
	Enforce        bool
	ExtractionKits []extractor.ExtractorKit
}

// NewExtAuthZFilter creates a new ExtAuthZFilter
func NewExtAuthZFilter(config Config, c *client.OpenFgaClient, logger logger.Logger) *ExtAuthZFilter {
	return &ExtAuthZFilter{enforce: config.Enforce, client: c, extractionKits: config.ExtractionKits, logger: logger}
}

func (e ExtAuthZFilter) Register(server *grpc.Server) {
	envoy.RegisterAuthorizationServer(server, e)
}

// Check the access decision based on the incoming request
func (e ExtAuthZFilter) Check(ctx context.Context, req *envoy.CheckRequest) (response *envoy.CheckResponse, err error) {
	reqID := req.Attributes.GetRequest().GetHttp().GetHeaders()["x-request-id"]
	logger := e.logger
	if reqID != "" {
		logger = e.logger.With(zap.String("request_id", reqID))
	}

	res, err := e.check(ctx, req, logger)
	if e.enforce {
		if err != nil {
			logger.Error("Failed to check permissions", zap.Error(err))
			return nil, err
		}

		return res, nil
	} else {
		if err != nil {
			logger.Error("Failed to check permissions", zap.Error(err))
		}

		return allow, nil
	}
}

func (e ExtAuthZFilter) extract(ctx context.Context, req *envoy.CheckRequest) (*extractor.Check, error) {
	for _, es := range e.extractionKits {
		e.logger.Debug("Extracting values", zap.String("extractor", es.Name))
		check, err := es.Extract(ctx, req)
		if err == nil {
			return check, nil
		}

		if errors.Is(err, extractor.ErrValueNotFound) {
			e.logger.Debug("Extracing value not found", zap.String("extraction_kit", es.Name), zap.Error(err))
			continue
		}

		return nil, err
	}

	return nil, nil
}

// Check implements the Check method of the Authorization interface.
func (e ExtAuthZFilter) check(ctx context.Context, req *envoy.CheckRequest, logger logger.Logger) (response *envoy.CheckResponse, err error) {
	check, err := e.extract(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("extracting values from request: %w", err)
	}

	if check == nil {
		logger.Error("Failed to extract values from request")
		return deny(codes.InvalidArgument, "No extraction set found"), nil
	}

	body := client.ClientCheckRequest{
		User:     check.User,
		Relation: check.Relation,
		Object:   check.Object,
		Context:  &check.Context,
	}

	options := client.ClientCheckOptions{
		AuthorizationModelId: openfga.PtrString(e.modelID),
	}

	e.logger.Debug("Checking permissions", zap.String("user", check.User), zap.String("relation", check.Relation), zap.String("object", check.Object))
	data, err := e.client.Check(ctx).Body(body).Options(options).Execute()
	if err != nil {
		logger.Error("Failed to check permissions", zap.Error(err))
		return deny(codes.Internal, fmt.Sprintf("Error checking permissions: %v", err)), nil
	}

	if data.GetAllowed() {
		logger.Debug("Access granted", zap.String("resolution", data.GetResolution()))
		return allow, nil
	}

	logger.Debug("Access denied")
	return deny(codes.PermissionDenied, fmt.Sprintf("Access denied: %s", data.GetResolution())), nil
}
