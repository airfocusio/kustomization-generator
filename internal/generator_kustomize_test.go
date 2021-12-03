package internal

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorKustomize(t *testing.T) {
	viperInstance := viper.New()

	c1, err := LoadGenerator(*viperInstance, "./generator_kustomize_test.yaml")
	if assert.NoError(t, err) {
		c2 := KustomizeGenerator{
			Namespace: "namespace",
			Url:       "github.com/owner/repo/kustomize?ref=main",
			Args:      []string{"--reorder", "legacy"},
		}
		assert.Equal(t, c2, *c1)
	}
}
