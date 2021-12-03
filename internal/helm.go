package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

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

func retrieveHelmChartUrl(registry string, chart string, version string) (*string, error) {
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
	body, err := ioutil.ReadAll(resp.Body)
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
