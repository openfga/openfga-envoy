package authz

import (
	"context"

	envoy "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

// NoopFilter is a noop implementation of the Envoy AuthZ filter.
type NoopFilter struct{}

var _ envoy.AuthorizationServer = NoopFilter{}

func (e NoopFilter) Register(server *grpc.Server) {
	envoy.RegisterAuthorizationServer(server, e)
}

func (e NoopFilter) Check(ctx context.Context, req *envoy.CheckRequest) (response *envoy.CheckResponse, err error) {
	return allow, nil
}
