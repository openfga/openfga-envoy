package extractor

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestMockUnmarshal(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := &MockConfig{}
		err := mapstructure.Decode(nil, c)
		require.NoError(t, err)
		require.Empty(t, c.Value)
	})

	t.Run("success", func(t *testing.T) {
		c := &MockConfig{}
		err := mapstructure.Decode(map[string]any{"value": "subject:user_123"}, c)
		require.NoError(t, err)
		require.Equal(t, "subject:user_123", c.Value)
	})
}
