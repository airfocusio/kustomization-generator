package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorDownload(t *testing.T) {
	c1, err := LoadGenerator("./generator_download_test.yaml")
	if assert.NoError(t, err) {
		c2 := DownloadGenerator{
			Url: "https://domain.com/resources.yaml",
		}
		assert.Equal(t, c2, *c1)
	}
}
