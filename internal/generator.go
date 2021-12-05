package internal

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

type Kustomization struct {
	Namespace string   `yaml:"namespace"`
	Resources []string `yaml:"resources"`
}

type Generator interface {
	Generate(dir string) (*Kustomization, error)
}

func LoadGenerator(viperInst viper.Viper, path string) (*Generator, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	err = readYaml(bytes, &raw)
	if err != nil {
		return nil, err
	}
	t, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("config is missing proper type")
	}

	var result Generator
	if t == "download" {
		generator := DownloadGenerator{}
		err = readYaml(bytes, &generator)
		if err != nil {
			return nil, err
		}
		result = generator
	}
	if t == "helm" {
		generator := HelmGenerator{}
		err = readYaml(bytes, &generator)
		if err != nil {
			return nil, err
		}
		result = generator
	}
	if t == "kustomize" {
		generator := KustomizeGenerator{}
		err = readYaml(bytes, &generator)
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
