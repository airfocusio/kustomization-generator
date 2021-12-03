package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

type DownloadGenerator struct {
	Namespace string `yaml:"namespace"`
	Url       string `yaml:"url"`
}

func (g DownloadGenerator) Generate(dir string) (*Kustomization, error) {
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

	kustomization := Kustomization{
		Namespace: g.Namespace,
		Resources: []string{"resources.yaml"},
	}
	err = ioutil.WriteFile(path.Join(dir, "resources.yaml"), body, 0o644)
	if err != nil {
		return nil, err
	}

	return &kustomization, nil
}
