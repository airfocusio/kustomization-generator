package internal

import (
	"fmt"
	"os/exec"
)

type KustomizeGenerator struct {
	Namespace string   `yaml:"namespace"`
	Url       string   `yaml:"url"`
	Args      []string `yaml:"args"`
}

func (g KustomizeGenerator) Generate() (*KustomizationWithEmbeddedResources, error) {
	kustomizePath, err := exec.LookPath("kustomize")
	if err != nil {
		return nil, fmt.Errorf("executing kustomize failed: executable not found")
	}
	kustomizeArgs := []string{
		"build",
		g.Url,
	}
	kustomizeArgs = append(kustomizeArgs, g.Args...)
	kustomizeStdout, kustomizeStderr, err := runCommand(*exec.Command(kustomizePath, kustomizeArgs...))
	if err != nil {
		return nil, fmt.Errorf("executing kustomize failed: %v\n%s", err, string(kustomizeStderr))
	}

	resources, err := splitCombinedKubernetesResources(string(kustomizeStdout))
	if err != nil {
		return nil, fmt.Errorf("splitting helm resources failed: %v", err)
	}
	result := KustomizationWithEmbeddedResources{
		Namespace: g.Namespace,
		Resources: resources,
	}
	return &result, nil
}
