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

func (g KustomizeGenerator) Generate(dir string) (*Kustomization, error) {
	kustomizePath, err := exec.LookPath("kustomize")
	if err != nil {
		return nil, fmt.Errorf("executing kustomize failed: executable not found")
	}
	kustomizeArgs := []string{
		"build",
		g.Url,
	}
	kustomizeArgs = append(kustomizeArgs, g.Args...)
	kustomizeOutput, err := exec.Command(kustomizePath, kustomizeArgs...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing kustomize failed: %v\n%s", err, string(kustomizeOutput))
	}

	kustomization := Kustomization{
		Namespace: g.Namespace,
		Resources: []string{"resources.yaml"},
	}
	err = ioutil.WriteFile(path.Join(dir, "resources.yaml"), kustomizeOutput, 0o644)
	if err != nil {
		return nil, err
	}

	return &kustomization, nil
}
