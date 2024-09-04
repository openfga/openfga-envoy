package extractor

import (
	"context"
	"errors"
	"strings"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"gopkg.in/yaml.v3"
)

const (
	clientCertHeader       = "x-forwarded-client-cert"
	spiffeKey              = "By="
	spiffeCurrentClientKey = "URI="
)

type SpiffeExtractionType int8

func (t SpiffeExtractionType) String() string {
	switch t {
	case SpiffeTypeUser:
		return "user"
	case SpiffeTypeObject:
		return "object"
	}

	return "unknown"
}

func (t SpiffeExtractionType) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *SpiffeExtractionType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "user":
		*t = SpiffeTypeUser
	case "object":
		*t = SpiffeTypeObject
	default:
		return errors.New("unknown spiffe extraction type")
	}

	return nil
}

const (
	SpiffeTypeUser SpiffeExtractionType = iota
	SpiffeTypeObject
)

type SpiffeConfig struct {
	Type SpiffeExtractionType `yaml:"type"`
}

func NewSpiffe(config *SpiffeConfig) Extractor {
	if config == nil {
		config = &SpiffeConfig{}
	}

	var prefix string

	if config.Type == SpiffeTypeUser {
		prefix = spiffeCurrentClientKey
	} else {
		prefix = spiffeKey
	}

	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		headers := value.GetAttributes().GetRequest().GetHttp().GetHeaders()
		val, ok := headers[clientCertHeader]
		if !ok {
			return Extraction{}, false, nil
		}

		var segments = strings.Split(val, ",")

		for _, seg := range segments {
			parts := strings.Split(seg, ";")
			for _, part := range parts {
				if !strings.HasPrefix(part, prefix) {
					continue
				}

				if part[len(prefix):] == "" {
					continue
				}

				return Extraction{Value: part[len(prefix):]}, true, nil
			}
		}

		return Extraction{}, false, nil
	}
}
