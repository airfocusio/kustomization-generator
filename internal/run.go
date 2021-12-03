package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func Run(dir string, config Generator) error {
	tempDir, err := ioutil.TempDir("", ".kustomization-generator-")
	if err != nil {
		return fmt.Errorf("preparing temporary folder failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	kustomization, err := config.Generate(tempDir)
	if err != nil {
		return err
	}
	for i, res := range kustomization.Resources {
		kustomization.Resources[i] = path.Join("generated", res)
	}
	err = writeYamlFile(path.Join(dir, "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}

	os.RemoveAll(path.Join(dir, "generated"))
	if err != nil {
		return fmt.Errorf("clearing generated files failed: %v", err)
	}
	err = copyDir(tempDir, path.Join(dir, "generated"))
	if err != nil {
		return fmt.Errorf("copying generated files failed: %v", err)
	}

	return nil
}
