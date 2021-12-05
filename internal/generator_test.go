package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCombinedKubernetesResources(t *testing.T) {
	var str string
	var err error
	var actual []KustomizationResource

	str = `
apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`
	actual, err = splitCombinedKubernetesResources(str)
	if assert.NoError(t, err) {
		expected := []KustomizationResource{
			{
				Name: "database-secret.yaml",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`,
			},
		}
		assert.Equal(t, expected, actual)
	}

	str = `
apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure

---

apiVersion: v1
kind: Secret
metadata:
  name: other
stringData:
  password: secure
`
	actual, err = splitCombinedKubernetesResources(str)
	if assert.NoError(t, err) {
		expected := []KustomizationResource{
			{
				Name: "database-secret.yaml",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`,
			},
			{
				Name: "other-secret.yaml",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: other
stringData:
  password: secure
`,
			},
		}
		assert.Equal(t, expected, actual)
	}

	str = `
apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure

---

apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`
	actual, err = splitCombinedKubernetesResources(str)
	if assert.NoError(t, err) {
		expected := []KustomizationResource{
			{
				Name: "database-secret.yaml",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`,
			},
			{
				Name: "database-secret-1.yaml",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`,
			},
		}
		assert.Equal(t, expected, actual)
	}
}

func TestGetUniqueKubernetesResourceFileName(t *testing.T) {
	state := map[string]int{}
	assert.Equal(t, "foo", getUniqueKubernetesResourceFileName("foo", &state))
	assert.Equal(t, "foo-1", getUniqueKubernetesResourceFileName("foo", &state))
	assert.Equal(t, "foo-2", getUniqueKubernetesResourceFileName("foo", &state))

	assert.Equal(t, "bar", getUniqueKubernetesResourceFileName("bar", &state))
	assert.Equal(t, "bar-1", getUniqueKubernetesResourceFileName("bar-1", &state))
	assert.Equal(t, "bar-1-1", getUniqueKubernetesResourceFileName("bar", &state))
	assert.Equal(t, "bar-2", getUniqueKubernetesResourceFileName("bar", &state))
}
