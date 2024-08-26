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
		require.Empty(t, c.User)
		require.Empty(t, c.Object)
		require.Empty(t, c.Relation)
		require.False(t, c.found)
	})

	t.Run("success", func(t *testing.T) {
		c := &MockConfig{}
		err := yaml.Unmarshal([]byte(`user: "subject:user_123"`), c)
		require.NoError(t, err)
		require.Equal(t, "subject:user_123", c.User)
		require.Empty(t, c.Object)
		require.Empty(t, c.Relation)
		require.True(t, c.found)
	})

	t.Run("fails because subject and object are passed", func(t *testing.T) {
		c := &MockConfig{}
		err := yaml.Unmarshal([]byte(`
user: "subject:user_123"
object: resource:service_abc`), c)
		require.Error(t, err)
	})
}
