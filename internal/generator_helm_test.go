package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorHelm(t *testing.T) {
	c1, err := LoadGenerator("./generator_helm_test.yaml")
	if assert.NoError(t, err) {
		c2 := HelmGenerator{
			Registry: "https://charts.domain.com",
			Chart:    "chart",
			Version:  "1.2.3",
			Args: []string{
				"--include-crds",
			},
			Values: map[string]interface{}{
				"foo": "bar",
			},
			Name:      "name",
			Namespace: "namespace",
		}
		assert.Equal(t, c2, *c1)
	}
}
