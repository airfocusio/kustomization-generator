package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

type Generator interface {
	Generate(dir string) error
}

func LoadGenerator(viperInst viper.Viper, path string) (*Generator, error) {
	var raw map[string]interface{}
	err := readYamlFile(path, &raw)
	if err != nil {
		return nil, err
	}
	t, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("config is missing proper type")
	}

	var result Generator
	if t == "helm" {
		generator := HelmGenerator{}
		err = readYamlFile(path, &generator)
		if err != nil {
			return nil, err
		}
		result = generator
	}
	if t == "kustomize" {
		generator := KustomizeGenerator{}
		err = readYamlFile(path, &generator)
		if err != nil {
			return nil, err
		}
		result = generator
	}

	if result == nil {
		return nil, fmt.Errorf("config has unknown type %s", t)
	}
	return &result, nil
}

type kustomization struct {
	Namespace string   `yaml:"namespace"`
	Resources []string `yaml:"resources"`
}
