package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func Run(dir string, config Generator) error {
	kustomizationWithEmbeddedResources, err := config.Generate()
	if err != nil {
		return err
	}

	cleanKustomization(dir)
	err = writeKustomization(dir, *kustomizationWithEmbeddedResources)
	if err != nil {
		return err
	}

	return nil
}

func cleanKustomization(dir string) {
	os.RemoveAll(path.Join(dir, "generated"))
	os.Remove(path.Join(dir, "kustomization.yaml"))
}

func writeKustomization(dir string, kustomization KustomizationWithEmbeddedResources) error {
	kustomizationRaw := Kustomization{
		Namespace: kustomization.Namespace,
		Resources: []string{},
	}
	for _, resource := range kustomization.Resources {
		kustomizationRaw.Resources = append(kustomizationRaw.Resources, path.Join("generated", resource.Name))
	}

	err := writeYamlFile(path.Join(dir, "kustomization.yaml"), kustomizationRaw)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}

	for _, resource := range kustomization.Resources {
		resPath := path.Join(dir, "generated", resource.Name)
		err := os.MkdirAll(path.Dir(resPath), 0o755)
		if err != nil {
			return fmt.Errorf("writing kustomization failed: %v", err)
		}
		err = ioutil.WriteFile(resPath, []byte(resource.Content), 0o644)
		if err != nil {
			return fmt.Errorf("writing kustomization resource failed: %v", err)
		}
	}

	return nil
}
