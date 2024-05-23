package extractor

import (
	"context"
	"errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// Extractor is the interface for extracting values from a CheckRequest.
type Extractor func(ctx context.Context, value *authv3.CheckRequest) (string, bool, error)

type ExtractorSet struct {
	Name     string
	Object   Extractor
	User     Extractor
	Relation Extractor
}

type Config interface{}

func GetExtractorConfig(name string) (Config, error) {
	switch name {
	case "mock":
		return &MockConfig{}, nil
	default:
		return nil, errors.New("extractor not found")
	}
}

func MakeExtractor(name string, cfg Config) (Extractor, error) {
	switch name {
	case "mock":
		return NewMock(cfg.(*MockConfig)), nil
	default:
		return nil, errors.New("extractor not found")
	}
}
