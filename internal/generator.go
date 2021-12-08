package internal

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/airfocusio/go-expandenv"
	"gopkg.in/yaml.v3"
)

type Kustomization struct {
	Resources []string `yaml:"resources"`
}

type GeneratorResource struct {
	ApiVersion string
	Kind       string
	File       string
	Content    string
}

type GeneratorResult struct {
	Resources []GeneratorResource
}

type Generator interface {
	Generate() (*GeneratorResult, error)
}

type KubernetesResourceMetadata struct {
	Name string `yaml:"name"`
}

type KubernetesResource struct {
	ApiVersion string                     `yaml:"apiVersion"`
	Kind       string                     `yaml:"kind"`
	Metadata   KubernetesResourceMetadata `yaml:"metadata"`
}

func (r KubernetesResource) NonEmpty() bool {
	return r.ApiVersion != "" && r.Kind != "" && r.Metadata.Name != ""
}

func LoadGenerator(file string) (*Generator, error) {
	bytesRaw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var expansionTemp interface{}
	err = yaml.Unmarshal(bytesRaw, &expansionTemp)
	if err != nil {
		return nil, err
	}
	expansionTemp, err = expandenv.ExpandEnv(expansionTemp)
	if err != nil {
		return nil, err
	}
	bytes, err := yaml.Marshal(expansionTemp)
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

func splitCombinedKubernetesResources(all string) ([]GeneratorResource, error) {
	newLine := "\n"
	seperator := "---"

	allLines := append(strings.Split(strings.ReplaceAll(strings.ReplaceAll(all, "\r\n", newLine), "\r", newLine), newLine), seperator)
	for i := range allLines {
		allLines[i] = strings.TrimRight(allLines[i], " \t")
	}
	result := []GeneratorResource{}
	existingNames := map[string]int{}

	start := 0
	for i, line := range allLines {
		empty := true
		for j := start; j < i; j++ {
			if allLines[j] != "---" && allLines[j] != "" && !strings.HasPrefix(strings.TrimLeft(allLines[j], " \t"), "#") {
				empty = false
				break
			}
		}
		if !empty && strings.HasPrefix(line, seperator) {
			lines := []string{}
			for _, line := range allLines[start:i] {
				if line != seperator {
					lines = append(lines, line)
				}
			}
			content := strings.Trim(strings.Join(lines, newLine), "\n \t") + newLine
			start = i + 1

			kubernetesResource := KubernetesResource{}
			err := yaml.Unmarshal([]byte(content), &kubernetesResource)
			if err != nil {
				return result, err
			}
			if !kubernetesResource.NonEmpty() {
				continue
			}

			nameBase := strings.Trim(fmt.Sprintf("%s-%s", kubernetesResource.Metadata.Name, kubernetesResource.Kind), "-")
			name := getUniqueKubernetesResourceFileName(nameBase, &existingNames)
			result = append(result, GeneratorResource{
				ApiVersion: kubernetesResource.ApiVersion,
				Kind:       kubernetesResource.Kind,
				File:       name + ".yaml",
				Content:    content,
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
