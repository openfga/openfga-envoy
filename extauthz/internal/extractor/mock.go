package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type MockConfig struct {
	Value   string                 `yaml:"value"`
	Context map[string]interface{} `yaml:"context"`
	Err     error                  `yaml:"error"`
}

func NewMock(cfg *MockConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return Extraction{
			Value:   cfg.Value,
			Context: cfg.Context,
		}, cfg.Value != "", cfg.Err
	}
}
