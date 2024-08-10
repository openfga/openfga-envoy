package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type MockConfig struct {
	Val     string                 `yaml:"value"`
	Context map[string]interface{} `yaml:"context"`
	Err     error                  `yaml:"error"`
}

func NewMock(cfg *MockConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return Extraction{Value: cfg.Val, Context: cfg.Context}, cfg.Val != "", cfg.Err
	}
}
