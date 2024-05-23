package extractor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type MockConfig struct {
	Val string `yaml:"value"`
	Err error  `yaml:"error"`
}

func NewMock(cfg *MockConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (string, bool, error) {
		return cfg.Val, cfg.Val != "", cfg.Err
	}
}
