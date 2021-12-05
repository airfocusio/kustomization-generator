package internal

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Kustomization struct {
	Namespace string   `yaml:"namespace"`
	Resources []string `yaml:"resources"`
}

type KustomizationResource struct {
	Name    string
	Content string
}

type KustomizationWithEmbeddedResources struct {
	Namespace string
	Resources []KustomizationResource
}

type Generator interface {
	Generate() (*KustomizationWithEmbeddedResources, error)
}

type KubernetesResourceMetadata struct {
	Name string `yaml:"name"`
}

type KubernetesResource struct {
	ApiVersion string                     `yaml:"apiVersion"`
	Kind       string                     `yaml:"kind"`
	Metadata   KubernetesResourceMetadata `yaml:"metadata"`
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

func splitCombinedKubernetesResources(all string) ([]KustomizationResource, error) {
	newLine := "\n"
	seperator := "---"

	allLines := append(strings.Split(strings.ReplaceAll(strings.ReplaceAll(all, "\r\n", newLine), "\r", newLine), newLine), seperator)
	for i := range allLines {
		allLines[i] = strings.TrimRight(allLines[i], " \t")
	}
	result := []KustomizationResource{}
	existingNames := map[string]int{}

	start := 0
	for i, line := range allLines {
		if strings.HasPrefix(line, seperator) {
			content := strings.Trim(strings.Join(allLines[start:i], newLine), "\n \t") + newLine
			start = i + 1

			if content == "\n" {
				continue
			}

			kubernetesResource := KubernetesResource{}
			err := yaml.Unmarshal([]byte(content), &kubernetesResource)
			if err != nil {
				return result, err
			}
			nameBase := strings.Trim(fmt.Sprintf("%s-%s", kubernetesResource.Metadata.Name, kubernetesResource.Kind), "-")
			if nameBase == "" {
				nameBase = "unnamed"
			}
			name := getUniqueKubernetesResourceFileName(nameBase, &existingNames)
			result = append(result, KustomizationResource{
				Name:    name + ".yaml",
				Content: content,
			})
		}
	}

	return result, nil
}

func getUniqueKubernetesResourceFileName(name string, existing *map[string]int) string {
	invalidRegex := regexp.MustCompile("[^a-z0-9]+")
	nameNormalized := invalidRegex.ReplaceAllString(strings.ToLower(name), "-")
	counter := (*existing)[nameNormalized]
	(*existing)[nameNormalized] = counter + 1
	if counter > 0 {
		return getUniqueKubernetesResourceFileName(fmt.Sprintf("%s-%d", nameNormalized, counter), existing)
	}
	return nameNormalized
}
