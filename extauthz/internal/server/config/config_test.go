package config

import (
	"testing"

	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := LoadConfig("testdata/config.yaml")
	require.NoError(t, err)
	require.Equal(t, "http://localhost:8080", cfg.Server.APIURL)
	require.Equal(t, "01FQH7V8BEG3GPQW93KTRFR8JB", cfg.Server.StoreID)
	require.Equal(t, "01GXSA8YR785C4FYS3C0RTG7B1", cfg.Server.AuthorizationModelID)

	require.Equal(t, "debug", cfg.Log.Level)
	require.Equal(t, "text", cfg.Log.Format)
	require.Equal(t, "ISO8601", cfg.Log.TimestampFormat)

	require.Len(t, cfg.ExtractionSet, 1)
	require.Equal(t, "test", cfg.ExtractionSet[0].Name)
	require.Equal(t, "mock", cfg.ExtractionSet[0].User.Type)
	require.Equal(t, "subject:my_user", cfg.ExtractionSet[0].User.Config.(*extractor.MockConfig).Value)
	require.Equal(t, "mock", cfg.ExtractionSet[0].Object.Type)
	require.Equal(t, "resource:my_object", cfg.ExtractionSet[0].Object.Config.(*extractor.MockConfig).Value)
	require.Equal(t, "mock", cfg.ExtractionSet[0].Relation.Type)
	require.Equal(t, "my_relation", cfg.ExtractionSet[0].Relation.Config.(*extractor.MockConfig).Value)
}
