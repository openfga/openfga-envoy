package extractor

import (
	"context"
	"errors"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type Check struct {
	User     string
	Relation string
	Object   string
	Context  map[string]interface{}
}

type Extraction func(*Check) error

// Extractor is the interface for extracting values from a CheckRequest.
type Extractor func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error)

type ExtractorKit struct {
	Name     string
	User     Extractor
	Object   Extractor
	Relation Extractor
}

var ErrValueNotFound = errors.New("extraction value not found")

func (ek ExtractorKit) Extract(ctx context.Context, req *authv3.CheckRequest) (*Check, error) {
	check := &Check{}

	eUser, found, err := ek.User(ctx, req)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("extracting user: %w", ErrValueNotFound)
	}

	if err := eUser(check); err != nil {
		return nil, err
	}

	eObject, found, err := ek.Object(ctx, req)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("extracting object: %w", ErrValueNotFound)
	}

	if err := eObject(check); err != nil {
		return nil, err
	}

	eRelation, found, err := ek.Relation(ctx, req)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, fmt.Errorf("extracting relation: %w", ErrValueNotFound)
	}

	if err := eRelation(check); err != nil {
		return nil, err
	}

	return check, nil
}

type Config interface{}

func GetExtractorConfig(name string) (Config, error) {
	switch name {
	case "mock":
		return &MockConfig{}, nil
	case "method":
		return nil, nil
	default:
		return nil, errors.New("extractor not found")
	}
}

func MakeExtractor(name string, cfg Config) (Extractor, error) {
	switch name {
	case "mock":
		return NewMock(cfg.(*MockConfig)), nil
	case "method":
		return NewMethod(cfg), nil
	default:
		return nil, errors.New("extractor not found")
	}
}
