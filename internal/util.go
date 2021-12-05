package internal

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func fileList(dir string, includes []regexp.Regexp, excludes []regexp.Regexp) (*[]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pathRel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		pathRel = "/" + pathRel

		for _, i := range includes {
			if i.Match([]byte(pathRel)) {
				excluded := false
				for _, e := range excludes {
					if e.MatchString(pathRel) {
						excluded = true
					}
				}
				if !excluded {
					files = append(files, path)
					return nil
				}
			}
		}
		return nil
	})
	return &files, err
}

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
	bytes, err := ioutil.ReadFile(file)
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
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src string, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func runCommand(cmd exec.Cmd) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
