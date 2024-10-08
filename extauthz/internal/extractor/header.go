package extractor

import (
	"context"
	"errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/mitchellh/mapstructure"
)

type HeaderConfig struct {
	Name string `mapstructure:"name"`
}

func (c *HeaderConfig) UnmarshalMap(data any) error {
	var rawConfig = struct {
		Name string `mapstructure:"name"`
	}{}

	if err := mapstructure.Decode(data, &rawConfig); err != nil {
		return err
	}

	if rawConfig.Name == "" {
		return errors.New("header name is required")
	}

	c.Name = rawConfig.Name

	return nil
}

func NewHeader(config *HeaderConfig) Extractor {
	if config == nil {
		config = &HeaderConfig{}
	}

	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		headers := value.GetAttributes().GetRequest().GetHttp().GetHeaders()
		val, ok := headers[config.Name]
		if !ok || val == "" {
			return Extraction{}, false, nil
		}

		return Extraction{Value: val}, true, nil
	}
}
