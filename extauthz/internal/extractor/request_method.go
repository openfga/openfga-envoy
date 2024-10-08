package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

func NewRequestMethod(any) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return Extraction{
			Value: "access",
			Context: map[string]any{
				"request_method": value.GetAttributes().GetRequest().GetHttp().GetMethod(),
			},
		}, true, nil
	}
}
