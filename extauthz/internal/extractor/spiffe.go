package extractor

import (
	"context"
	"errors"
	"strings"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
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

func (t *spiffeExtractionType) UnmarshalMap(value any) error {
	switch value.(string) {
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
	Type spiffeExtractionType `mapstructure:"type"`
}

func NewSpiffe(config *SpiffeConfig) Extractor {
	if config == nil {
		config = &SpiffeConfig{}
	}

	var prefix string

	if config.Type == spiffeTypeUser {
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
