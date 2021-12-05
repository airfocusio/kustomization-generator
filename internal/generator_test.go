package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCombinedKubernetesResources(t *testing.T) {
	str := `apiVersion: v1
kind: Secret
metadata:
  name: database
stringData:
  password: secure
`
	actual, err := splitCombinedKubernetesResources([]byte(str))
	if assert.NoError(t, err) {
		expected := map[string][]byte{
			"resources.yaml": []byte(str),
		}
		assert.Equal(t, expected, *actual)
	}
}
