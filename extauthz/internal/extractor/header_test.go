package extractor

import (
	"context"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestHeaderUnmarshal(t *testing.T) {
	c := &HeaderConfig{}
	err := mapstructure.Decode("", c)
	require.Error(t, err)

	c = &HeaderConfig{}
	err = mapstructure.Decode(map[string]any{"name": "x-user-id"}, c)
	require.NoError(t, err)
	require.Equal(t, "x-user-id", c.Name)
}

func TestHeaderExtractor(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		extractor := NewHeader(nil)

		extraction, found, err := extractor(context.Background(), &authv3.CheckRequest{
			Attributes: &authv3.AttributeContext{
				Request: &authv3.AttributeContext_Request{
					Http: &authv3.AttributeContext_HttpRequest{},
				},
			},
		})

		require.NoError(t, err)
		require.False(t, found)
		require.Empty(t, extraction.Value)
	})

	t.Run("success for subject", func(t *testing.T) {
		extractor := NewHeader(&HeaderConfig{
			Name: "x-user-id",
		})

		extraction, found, err := extractor(context.Background(), &authv3.CheckRequest{
			Attributes: &authv3.AttributeContext{
				Request: &authv3.AttributeContext_Request{
					Http: &authv3.AttributeContext_HttpRequest{
						Headers: map[string]string{
							"x-user-id": "alice",
						},
					},
				},
			},
		})

		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, "alice", extraction.Value)
	})
}
