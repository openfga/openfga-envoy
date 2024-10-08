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
	Context  map[string]any
}

func (c Check) validate() error {
	if c.User == "" {
		return errors.New("user is required")
	}

	if c.Object == "" {
		return errors.New("object is required")
	}

	if c.Relation == "" {
		return errors.New("relation is required")
	}

	return nil
}

type Extraction struct {
	Value   string
	Context map[string]any
}

func (e *Extraction) applyExtraction(v *string, context map[string]any) error {
	*v = e.Value
	for k, v := range e.Context {
		if _, ok := context[k]; ok {
			return fmt.Errorf("context key %s already exists", k)
		}
		context[k] = v
	}
	return nil
}

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
	check := &Check{
		Context: make(map[string]any),
	}

	eUser, found, err := ek.User(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting user extraction: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("getting user extraction: %w", ErrValueNotFound)
	}

	if err := eUser.applyExtraction(&check.User, check.Context); err != nil {
		return nil, fmt.Errorf("extracting user: %w", err)
	}

	eObject, found, err := ek.Object(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting object extraction: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("getting object extraction: %w", ErrValueNotFound)
	}

	if err := eObject.applyExtraction(&check.Object, check.Context); err != nil {
		return nil, fmt.Errorf("extracting object: %w", err)
	}

	eRelation, found, err := ek.Relation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting relation extraction: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("getting relation extraction: %w", ErrValueNotFound)
	}

	if err := eRelation.applyExtraction(&check.Relation, check.Context); err != nil {
		return nil, fmt.Errorf("extracting relation: %w", err)
	}

	if err := check.validate(); err != nil {
		return nil, fmt.Errorf("validating check: %v", err)
	}

	return check, nil
}

type Config any

func GetExtractorConfig(name string) (Config, error) {
	switch name {
	case "mock":
		return &MockConfig{}, nil
	case "header":
		return &HeaderConfig{}, nil
	case "request_method":
		return nil, nil
	case "spiffe":
		return &SpiffeConfig{}, nil
	default:
		return nil, errors.New("extractor not found")
	}
}

func MakeExtractor(name string, cfg Config) (Extractor, error) {
	switch name {
	case "mock":
		return NewMock(cfg.(*MockConfig)), nil
	case "header":
		return NewHeader(cfg.(*HeaderConfig)), nil
	case "request_method":
		return NewRequestMethod(cfg), nil
	case "spiffe":
		return NewSpiffe(cfg.(*SpiffeConfig)), nil
	default:
		return nil, errors.New("extractor not found")
	}
}
