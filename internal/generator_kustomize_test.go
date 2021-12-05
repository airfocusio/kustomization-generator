package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorKustomize(t *testing.T) {
	c1, err := LoadGenerator("./generator_kustomize_test.yaml")
	if assert.NoError(t, err) {
		c2 := KustomizeGenerator{
			Namespace: "namespace",
			Url:       "github.com/owner/repo/kustomize?ref=main",
			Args:      []string{"--reorder", "legacy"},
		}
		assert.Equal(t, c2, *c1)
	}
}
