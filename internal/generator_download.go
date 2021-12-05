package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type DownloadGenerator struct {
	Namespace string `yaml:"namespace"`
	Url       string `yaml:"url"`
}

func (g DownloadGenerator) Generate() (*KustomizationWithEmbeddedResources, error) {
	req, err := http.NewRequest("GET", g.Url, nil)
	client := &http.Client{}
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %v", g.Url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %v", g.Url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %v", g.Url, err)
	}

	resources, err := splitCombinedKubernetesResources(string(body))
	if err != nil {
		return nil, fmt.Errorf("splitting helm resources failed: %v", err)
	}
	result := KustomizationWithEmbeddedResources{
		Namespace: g.Namespace,
		Resources: resources,
	}
	return &result, nil
}
