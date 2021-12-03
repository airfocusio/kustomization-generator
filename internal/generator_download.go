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

func (g DownloadGenerator) Generate(dir string) error {
	req, err := http.NewRequest("GET", g.Url, nil)
	client := &http.Client{}
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", g.Url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", g.Url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", g.Url, err)
	}

	kustomization := kustomization{
		Namespace: g.Namespace,
		Resources: []string{"resources.yaml"},
	}
	err = writeYamlFile(path.Join(dir, "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}
	err = ioutil.WriteFile(path.Join(dir, "resources.yaml"), body, 0o644)
	if err != nil {
		return fmt.Errorf("writing generated resource failed: %v", err)
	}

	return nil
}
