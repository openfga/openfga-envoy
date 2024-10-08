package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/openfga/openfga-envoy/extauthz/internal/extractor"
	"github.com/spf13/viper"
)

type Server struct {
	APIURL               string `mapstructure:"api_url"`
	StoreID              string `mapstructure:"store_id"`
	AuthorizationModelID string `mapstructure:"authorization_model_id"`
}

type Extractor struct {
	Type   string           `mapstructure:"type"`
	Config extractor.Config `mapstructure:"config"`
}

func (e *Extractor) UnmarshalMap(data any) error {
	var err error
	var rawConfig struct {
		Type   string         `mapstructure:"type"`
		Config map[string]any `mapstructure:"config"`
	}

	if err = mapstructure.Decode(data, &rawConfig); err != nil {
		return err
	}

	e.Type = rawConfig.Type
	e.Config, err = extractor.GetExtractorConfig(rawConfig.Type)
	if err != nil {
		return fmt.Errorf("getting %s: %w", e.Type, err)
	}

	if e.Config != nil {
		if err := mapstructure.Decode(rawConfig.Config, e.Config); err != nil {
			return err
		}
	}

	return nil
}

type ExtractionSet struct {
	Name     string    `mapstructure:"name"`
	User     Extractor `mapstructure:"user"`
	Object   Extractor `mapstructure:"object"`
	Relation Extractor `mapstructure:"relation"`
}

type AuthMode int8

const (
	AuthModeMonitor AuthMode = iota + 1
	AuthModeEnforce
	AuthModeDisabled
)

func (m AuthMode) String() string {
	switch m {
	case AuthModeMonitor:
		return "MONITOR"
	case AuthModeEnforce:
		return "ENFORCE"
	case AuthModeDisabled:
		return "DISABLED"
	}

	return "UNKNOWN"
}

func (m *AuthMode) UnmarshalMap(value any) error {
	switch value.(string) {
	case "ENFORCE":
		*m = AuthModeEnforce
	case "DISABLED":
		*m = AuthModeDisabled
	case "MONITOR":
		*m = AuthModeMonitor
	default:
		return errors.New("unknown mode")
	}

	return nil
}

type Config struct {
	ExtractionSets []ExtractionSet `mapstructure:"extraction_sets"`
	Server         Server          `mapstructure:"server"`
	Log            Log             `mapstructure:"log"`
	Mode           AuthMode        `mapstructure:"mode"`
}

type Log struct {
	Level           string `mapstructure:"level"`
	Format          string `mapstructure:"format"`
	TimestampFormat string `mapstructure:"timestamp_format"`
}

func LoadConfig(path string) (Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := initDefaultValues(v); err != nil {
		return Config{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	v.SetEnvPrefix("OPENFGA_EXTAUTHZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if err := v.MergeConfig(f); err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := v.Unmarshal(
		&cfg,
		viper.DecodeHook(func(st reflect.Type, tt reflect.Type, data any) (any, error) {
			if tt == reflect.TypeOf(Extractor{}) {
				es := new(Extractor)
				if err := es.UnmarshalMap(data); err != nil {
					return nil, err
				}

				return es, nil
			}

			if tt == reflect.TypeOf(AuthMode(0)) {
				am := new(AuthMode)
				if err := am.UnmarshalMap(data); err != nil {
					return nil, err
				}

				return am, nil
			}

			return data, nil
		}),
	); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func initDefaultValues(v *viper.Viper) error {
	defaultConfig := Config{}
	defaultConfig.Log.Level = "info"
	defaultConfig.Log.Format = "json"
	defaultConfig.Log.TimestampFormat = "RFC3339"
	defaultConfig.Mode = AuthModeDisabled

	return loadConfigFromStruct(v, defaultConfig)
}

// loadConfigFromStruct loads the configuration from a struct allowing to
// load env variables for all fields despite whether they are included in the config
// file or not.
// see https://github.com/spf13/viper/issues/584#issuecomment-691141622
func loadConfigFromStruct(v *viper.Viper, cfg any) error {
	cfgMap := make(map[string]any)
	if err := mapstructure.Decode(cfg, &cfgMap); err != nil {
		return err
	}

	cfgJsonBytes, err := json.Marshal(&cfgMap)
	if err != nil {
		return err
	}

	return v.ReadConfig(bytes.NewReader(cfgJsonBytes))
}
