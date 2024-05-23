package config

import (
	"fmt"
	"os"

	"github.com/jcchavezs/openfga-envoy/extauthz/internal/extractor"
	"gopkg.in/yaml.v3"
)

type Server struct {
	APIURL               string `yaml:"api_url"`
	StoreID              string `yaml:"store_id"`
	AuthorizationModelID string `yaml:"authorization_model_id"`
}

type Extractor struct {
	Type   string           `yaml:"type"`
	Config extractor.Config `yaml:"config"`
}

func (e *Extractor) UnmarshalYAML(value *yaml.Node) error {
	var err error
	var rawConfig struct {
		Type   string    `yaml:"type"`
		Config yaml.Node `yaml:"config"`
	}

	if err = value.Decode(&rawConfig); err != nil {
		return err
	}

	e.Type = rawConfig.Type
	e.Config, err = extractor.GetExtractorConfig(rawConfig.Type)
	if err != nil {
		return fmt.Errorf("getting %s: %w", e.Type, err)
	}

	if err := rawConfig.Config.Decode(e.Config); err != nil {
		return err
	}

	return nil
}

type ExtractionSet struct {
	Name     string    `yaml:"name"`
	User     Extractor `yaml:"user"`
	Object   Extractor `yaml:"object"`
	Relation Extractor `yaml:"relation"`
}

type Config struct {
	ExtractionSet []ExtractionSet `yaml:"extraction_sets"`
	Server        Server          `yaml:"server"`
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	config, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(config, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
