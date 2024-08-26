package extractor

import (
	"context"
	"errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"gopkg.in/yaml.v3"
)

type MockConfig struct {
	found    bool
	User     string                 `yaml:"user"`
	Object   string                 `yaml:"object"`
	Relation string                 `yaml:"relation"`
	Context  map[string]interface{} `yaml:"context"`
	Err      error                  `yaml:"error"`
	CheckErr error                  `yaml:"check_error"`
}

const (
	hasUser = 1 << iota
	hasObject
	hasRelation
)

func (c *MockConfig) UnmarshalYAML(value *yaml.Node) error {
	type mockConfig MockConfig
	if err := value.Decode((*mockConfig)(c)); err != nil {
		return err
	}

	ec := 0
	if c.User != "" {
		ec |= hasUser
	}

	if c.Object != "" {
		ec |= hasObject
	}

	if c.Relation != "" {
		ec |= hasRelation
	}

	if ec == 0 {
		return nil
	}

	c.found = true

	if (ec & -ec) != ec {
		return errors.New("only one of user, object, or relation can be set")
	}

	return nil
}

func NewMock(cfg *MockConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		return func(c *Check) error {
			if cfg.CheckErr != nil {
				return cfg.CheckErr
			}

			c.User = cfg.User
			c.Object = cfg.Object
			c.Relation = cfg.Relation
			for k, v := range cfg.Context {
				c.Context[k] = v
			}

			return nil
		}, cfg.found, cfg.Err
	}
}
