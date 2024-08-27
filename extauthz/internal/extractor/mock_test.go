package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMockUnmarshal(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := &MockConfig{}
		err := yaml.Unmarshal(nil, c)
		require.NoError(t, err)
		require.Empty(t, c.Value)
	})

	t.Run("success", func(t *testing.T) {
		c := &MockConfig{}
		err := yaml.Unmarshal([]byte(`value: "subject:user_123"`), c)
		require.NoError(t, err)
		require.Equal(t, "subject:user_123", c.Value)
	})
}
