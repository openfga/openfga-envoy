package authz

import (
	"context"
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
	extractionSet []extractor.ExtractorSet
	modelID       string
}

var _ envoy.AuthorizationServer = (*ExtAuthZFilter)(nil)

// NewExtAuthZFilter creates a new ExtAuthZFilter
func NewExtAuthZFilter(c *client.OpenFgaClient, es []extractor.ExtractorSet) *ExtAuthZFilter {
	return &ExtAuthZFilter{client: c, extractionSet: es}
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

	log.Println(res)
	return res, nil
}

type extracted struct {
	user     extractor.Extraction
	object   extractor.Extraction
	relation extractor.Extraction
}

func (e ExtAuthZFilter) extract(ctx context.Context, req *envoy.CheckRequest) (*extracted, error) {
	var user, object, relation extractor.Extraction
	for _, es := range e.extractionSet {
		var (
			found bool
			err   error
		)
		user, found, err = es.User(ctx, req)
		if err != nil {
			return nil, err
		}

		if !found {
			continue
		}

		object, found, err = es.Object(ctx, req)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}

		relation, found, err = es.Relation(ctx, req)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}

		return &extracted{
			user:     user,
			object:   object,
			relation: relation,
		}, nil
	}

	return nil, nil
}

func mergeMaps(map1, map2 map[string]any) map[string]any {
	UniqueMap := make(map[string]any)

	// for loop for the first map
	for key, val := range map1 {
		UniqueMap[key] = val
	}

	// for loop for the second map
	for key, val := range map2 {
		UniqueMap[key] = val
	}
	// return merged result
	return UniqueMap
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

	context := map[string]any{}

	if extracted.user.Context != nil {
		context = mergeMaps(context, extracted.user.Context)
	}

	if extracted.object.Context != nil {
		context = mergeMaps(context, extracted.object.Context)
	}

	if extracted.relation.Context != nil {
		context = mergeMaps(context, extracted.relation.Context)
	}

	body := client.ClientCheckRequest{
		User:     extracted.user.Value,
		Relation: extracted.relation.Value,
		Object:   extracted.object.Value,
		Context:  &context,
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
