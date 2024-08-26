package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

func NewMethod(any) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return func(c *Check) error {
			c.Relation = "access"
			c.Context["method"] = value.GetAttributes().GetRequest().GetHttp().GetMethod()
			return nil
		}, true, nil
	}
}
