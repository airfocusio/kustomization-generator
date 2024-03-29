package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type HelmGenerator struct {
	Registry    string                 `yaml:"registry"`
	Chart       string                 `yaml:"chart"`
	Version     string                 `yaml:"version"`
	Name        string                 `yaml:"name"`
	Namespace   string                 `yaml:"namespace"`
	ApiVersions []string               `yaml:"apiVersions"`
	Args        []string               `yaml:"args"`
	Values      map[string]interface{} `yaml:"values"`
}

func (g HelmGenerator) Generate() (*GeneratorResult, error) {
	valuesPath, err := os.CreateTemp("", ".kustomization-generator-*-values.yaml")
	if err != nil {
		return nil, fmt.Errorf("writing temporary values file failed: %v", err)
	}
	defer os.Remove(valuesPath.Name())
	valuesBytes, err := yaml.Marshal(g.Values)
	if err != nil {
		return nil, fmt.Errorf("writing temporary values file failed: %v", err)
	}
	err = os.WriteFile(valuesPath.Name(), valuesBytes, 0o600)
	if err != nil {
		return nil, fmt.Errorf("writing temporary values file failed: %v", err)
	}

	helmPath, err := exec.LookPath("helm")
	if err != nil {
		return nil, fmt.Errorf("executing helm failed: executable not found")
	}
	helmArgs := []string{
		"template",
		g.Name,
		"--namespace", g.Namespace,
		"--values", valuesPath.Name(),
	}

	if strings.HasPrefix(g.Registry, "oci://") {
		helmArgs = append(helmArgs, g.Registry, "--version", g.Version)
	} else if strings.HasPrefix(g.Registry, "https://") {
		url, err := retrieveHelmChartArchiveUrl(g.Registry, g.Chart, g.Version)
		if err != nil {
			return nil, err
		}
		helmArgs = append(helmArgs, *url)
	} else {
		return nil, fmt.Errorf("unsupported registry %s", g.Registry)
	}

	if len(g.ApiVersions) > 0 {
		helmArgs = append(helmArgs, "--api-versions", strings.Join(g.ApiVersions, ","))
	}
	helmArgs = append(helmArgs, g.Args...)
	helmStdout, helmStderr, err := runCommand(*exec.Command(helmPath, helmArgs...))
	if err != nil {
		return nil, fmt.Errorf("executing helm failed: %v\n%s", err, string(helmStderr))
	}

	resources, err := splitCombinedKubernetesResources(string(helmStdout))
	if err != nil {
		return nil, fmt.Errorf("splitting helm resources failed: %v", err)
	}
	result := GeneratorResult{
		Resources: resources,
	}
	return &result, nil
}

type helmRegistryIndex struct {
	ApiVersion string `yaml:"apiVersion"`
	Entries    map[string][]struct {
		ApiVersion string   `yaml:"apiVersion"`
		AppVersion string   `yaml:"appVersion"`
		Name       string   `yaml:"name"`
		Version    string   `yaml:"version"`
		Urls       []string `yaml:"urls"`
	} `yaml:"entries"`
}

func retrieveHelmChartArchiveUrl(registry string, chart string, version string) (*string, error) {
	url := strings.TrimSuffix(registry, "/") + "/index.yaml"
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry index at %s: %v", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry index at %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry index at %s: %v", url, err)
	}
	index := helmRegistryIndex{}
	err = yaml.Unmarshal(body, &index)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry index at %s: %v", url, err)
	}

	versions, ok := index.Entries[chart]
	if !ok {
		return nil, fmt.Errorf("chart %s could not be found", chart)
	}
	for _, entry := range versions {
		if entry.Version == version {
			if len(entry.Urls) == 0 {
				return nil, fmt.Errorf("chart %s version %s has no download urls", chart, version)
			}
			if len(entry.Urls) > 1 {
				return nil, fmt.Errorf("chart %s version %s has multiple download urls", chart, version)
			}
			result := entry.Urls[0]
			if !strings.HasPrefix(result, "http://") && !strings.HasPrefix(result, "https://") {
				result = strings.TrimSuffix(registry, "/") + "/" + strings.TrimPrefix(result, "/")
			}
			return &result, nil
		}
	}
	return nil, fmt.Errorf("chart %s version %s could not be found", chart, version)
}
