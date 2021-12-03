package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func Run(dir string, config Config) error {
	url, err := retrieveHelmChartUrl(config.Registry, config.Chart, config.Version)
	if err != nil {
		return err
	}

	valuesPath, err := ioutil.TempFile("", ".tmp-*-values.yaml")
	if err != nil {
		return fmt.Errorf("writing temporary values file failed: %v", err)
	}
	defer os.Remove(valuesPath.Name())
	valuesBytes, err := yaml.Marshal(config.Values)
	if err != nil {
		return fmt.Errorf("writing temporary values file failed: %v", err)
	}
	err = os.WriteFile(valuesPath.Name(), valuesBytes, 0o600)
	if err != nil {
		return fmt.Errorf("writing temporary values file failed: %v", err)
	}

	tempDir, err := ioutil.TempDir("", ".tmp-")
	if err != nil {
		return fmt.Errorf("preparing temporary folder failed: %v", err)
	}
	defer os.RemoveAll(tempDir)
	helmPath, err := exec.LookPath("helm")
	if err != nil {
		return fmt.Errorf("executing helm failed: executable not found")
	}
	helmArgs := []string{
		"template",
		config.Name,
		*url,
		"--namespace", config.Namespace,
		"--output-dir", tempDir,
		"--values", valuesPath.Name(),
	}
	helmArgs = append(helmArgs, config.Args...)
	helmOutput, err := exec.Command(helmPath, helmArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("executing helm failed: %v\n%s", err, string(helmOutput))
	}

	kustomization := kustomization{
		Namespace: config.Namespace,
	}
	tempDir2 := path.Join(tempDir, config.Chart)
	includes := []regexp.Regexp{*regexp.MustCompile(`\.ya?ml$`)}
	excludes := []regexp.Regexp{}
	files, err := fileList(tempDir2, includes, excludes)
	if err != nil {
		return fmt.Errorf("listing helm generated resources failed: %v", err)
	}
	for _, file := range *files {
		rel, err := filepath.Rel(tempDir2, file)
		if err != nil {
			return fmt.Errorf("listing helm generated resources failed: %v", err)
		}
		kustomization.Bases = append(kustomization.Bases, rel)
	}

	err = os.RemoveAll(path.Join(dir, "crds"))
	if err != nil {
		return fmt.Errorf("cleaning up target failed: %v", err)
	}
	err = os.RemoveAll(path.Join(dir, "templates"))
	if err != nil {
		return fmt.Errorf("cleaning up target failed: %v", err)
	}

	err = writeYamlFile(path.Join(tempDir2, "kustomization.yaml"), kustomization)
	if err != nil {
		return fmt.Errorf("writing kustomization failed: %v", err)
	}
	err = copyDir(tempDir2, dir)
	if err != nil {
		return fmt.Errorf("copying files to target failed: %v", err)
	}

	return nil
}

type kustomization struct {
	Namespace string   `yaml:"namespace"`
	Bases     []string `yaml:"bases"`
}
