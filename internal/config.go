package internal

import (
	"github.com/spf13/viper"
)

type Config struct {
	Registry  string                 `yaml:"registry"`
	Chart     string                 `yaml:"chart"`
	Version   string                 `yaml:"version"`
	Name      string                 `yaml:"name"`
	Namespace string                 `yaml:"namespace"`
	Args      []string               `yaml:"args"`
	Values    map[string]interface{} `yaml:"values"`
}

func LoadConfig(viperInst viper.Viper, path string) (*Config, error) {
	config := Config{}
	err := readYamlFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
