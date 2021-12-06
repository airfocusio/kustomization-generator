package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const configFile = "kustomization-generator.yaml"

func Run(dir string) error {
	file := path.Join(dir, configFile)
	generator, err := LoadGenerator(file)
	if err != nil {
		return fmt.Errorf("unable to load configuration: %v", err)
	}

	kustomizationWithEmbeddedResources, err := (*generator).Generate()
	if err != nil {
		return err
	}

	clear(dir)
	err = write(dir, *kustomizationWithEmbeddedResources)
	if err != nil {
		return err
	}

	return nil
}

func clear(dir string) {
	fds, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fd := range fds {
		if fd.Name() == configFile {
			continue
		}

		fp := path.Join(dir, fd.Name())
		if fd.IsDir() {
			os.RemoveAll(fp)
		} else {
			os.Remove(fp)
		}
	}
}

func write(dir string, result GeneratorResult) error {
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

	kustomization := Kustomization{}

	for i := range buckets {
		bucket := &buckets[i]
		bucket.dir = path.Join(dir, bucket.name)
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
		if len(bucket.kustomization.Resources) > 0 {
			kustomization.Resources = append(kustomization.Resources, bucket.name)
			err := writeYamlFile(path.Join(bucket.dir, "kustomization.yaml"), bucket.kustomization)
			if err != nil {
				return fmt.Errorf("writing kustomization failed: %v", err)
			}
		}
	}

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}
	err = writeYamlFile(path.Join(dir, "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}

	return nil
}
