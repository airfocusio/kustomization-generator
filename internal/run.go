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
}

func writeKustomization(dir string, result GeneratorResult) error {
	buckets := []struct {
		name          string
		filter        func(resource GeneratorResource) bool
		kustomization Kustomization
		dir           string
	}{
		{
			name: "crds",
			filter: func(resource GeneratorResource) bool {
				return resource.ApiVersion == "apiextensions.k8s.io/v1" && resource.Kind == "CustomResourceDefinition"
			},
		},
		{
			name: "namespaces",
			filter: func(resource GeneratorResource) bool {
				return resource.ApiVersion == "v1" && resource.Kind == "Namespace"
			},
		},
		{
			name: "resources",
			filter: func(resource GeneratorResource) bool {
				return true
			},
		},
	}

	kustomization := Kustomization{
		Namespace: result.Namespace,
	}

	for i := range buckets {
		bucket := &buckets[i]
		bucket.kustomization.Namespace = result.Namespace
		bucket.dir = path.Join(dir, "generated", bucket.name)
		err := os.MkdirAll(bucket.dir, 0o755)
		if err != nil {
			return fmt.Errorf("writing kustomization failed: %v", err)
		}
	}

	for _, resource := range result.Resources {
		for i := range buckets {
			bucket := &buckets[i]
			if bucket.filter(resource) {
				bucket.kustomization.Resources = append(bucket.kustomization.Resources, resource.File)
				resourcePath := path.Join(bucket.dir, resource.File)
				err := ioutil.WriteFile(resourcePath, []byte(resource.Content), 0o644)
				if err != nil {
					return fmt.Errorf("writing kustomization failed: %v", err)
				}
				break
			}
		}
	}

	for _, bucket := range buckets {
		bucket.kustomization.Namespace = result.Namespace
		if len(bucket.kustomization.Resources) > 0 {
			kustomization.Resources = append(kustomization.Resources, bucket.name)
			err := writeYamlFile(path.Join(bucket.dir, "kustomization.yaml"), bucket.kustomization)
			if err != nil {
				return fmt.Errorf("writing kustomization failed: %v", err)
			}
		}
	}

	err := os.MkdirAll(path.Join(dir, "generated"), 0o755)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}
	err = writeYamlFile(path.Join(dir, "generated", "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}

	return nil
}
