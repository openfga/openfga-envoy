package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type MockConfig struct {
	Value   string         `mapstructure:"value"`
	Context map[string]any `mapstructure:"context"`
	Err     error          `mapstructure:"error"`
}

func NewMock(cfg *MockConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return Extraction{
			Value:   cfg.Value,
			Context: cfg.Context,
		}, cfg.Value != "", cfg.Err
	}
}
