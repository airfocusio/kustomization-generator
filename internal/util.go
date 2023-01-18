package internal

import (
	"bytes"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func readYaml(bytes []byte, v interface{}) error {
	return yaml.Unmarshal(bytes, v)
}

func writeYaml(v interface{}) ([]byte, error) {
	var bs bytes.Buffer
	enc := yaml.NewEncoder(&bs)
	enc.SetIndent(2)
	err := enc.Encode(v)
	return bs.Bytes(), err
}

func readYamlFile(file string, v interface{}) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = readYaml(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

func writeYamlFile(file string, v interface{}) error {
	bytes, err := writeYaml(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func runCommand(cmd exec.Cmd) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
