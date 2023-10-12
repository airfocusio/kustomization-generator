package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunDownload(t *testing.T) {
	err := Run("../example/download")
	assert.NoError(t, err)
}

func TestRunHelm(t *testing.T) {
	err := Run("../example/helm")
	assert.NoError(t, err)
}

func TestRunHelmOci(t *testing.T) {
	err := Run("../example/helm-oci")
	assert.NoError(t, err)
}

func TestRunKustomize(t *testing.T) {
	err := Run("../example/kustomize")
	assert.NoError(t, err)
}
