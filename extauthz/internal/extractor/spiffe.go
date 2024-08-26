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

type spiffeExtractionType int8

func (t spiffeExtractionType) String() string {
	switch t {
	case spiffeTypeUser:
		return "user"
	case spiffeTypeObject:
		return "object"
	}

	return "unknown"
}

func (t spiffeExtractionType) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *spiffeExtractionType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "user":
		*t = spiffeTypeUser
	case "object":
		*t = spiffeTypeObject
	default:
		return errors.New("unknown spiffe extraction type")
	}

	return nil
}

const (
	spiffeTypeUser spiffeExtractionType = iota
	spiffeTypeObject
)

type SpiffeConfig struct {
	_type spiffeExtractionType `yaml:"type"`
}

func NewSpiffe(config *SpiffeConfig) Extractor {
	return func(ctx context.Context, value *authv3.CheckRequest) (Extraction, bool, error) {
		headers := value.GetAttributes().GetRequest().GetHttp().GetHeaders()
		val, ok := headers[clientCertHeader]
		if !ok {
			return nil, false, nil
		}

		var prefix string
		if config._type == spiffeTypeUser {
			prefix = spiffeCurrentClientKey
		} else {
			prefix = spiffeKey
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

				if config._type == spiffeTypeUser {
					return func(c *Check) error {
						c.User = part[len(prefix):]
						return nil
					}, true, nil
				} else {
					return func(c *Check) error {
						c.Object = part[len(prefix):]
						return nil
					}, true, nil
				}
			}
		}

		return nil, false, nil
	}
}
