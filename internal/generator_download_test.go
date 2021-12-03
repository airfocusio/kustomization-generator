package internal

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorDownload(t *testing.T) {
	viperInstance := viper.New()

	c1, err := LoadGenerator(*viperInstance, "./generator_download_test.yaml")
	if assert.NoError(t, err) {
		c2 := DownloadGenerator{
			Namespace: "namespace",
			Url:       "https://domain.com/resources.yaml",
		}
		assert.Equal(t, c2, *c1)
	}
}
