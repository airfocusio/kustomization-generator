package internal

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
)

type KustomizeGenerator struct {
	Namespace string   `yaml:"namespace"`
	Url       string   `yaml:"url"`
	Args      []string `yaml:"args"`
}

func (g KustomizeGenerator) Generate(dir string) error {
	kustomizePath, err := exec.LookPath("kustomize")
	if err != nil {
		return fmt.Errorf("executing kustomize failed: executable not found")
	}
	kustomizeArgs := []string{
		"build",
		g.Url,
	}
	kustomizeArgs = append(kustomizeArgs, g.Args...)
	kustomizeOutput, err := exec.Command(kustomizePath, kustomizeArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("executing kustomize failed: %v\n%s", err, string(kustomizeOutput))
	}

	kustomization := kustomization{
		Namespace: g.Namespace,
		Resources: []string{"resources.yaml"},
	}
	err = writeYamlFile(path.Join(dir, "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}
	err = ioutil.WriteFile(path.Join(dir, "resources.yaml"), kustomizeOutput, 0o644)
	if err != nil {
		return fmt.Errorf("writing generated resource failed: %v", err)
	}

	return nil
}
