package internal

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	viperInstance := viper.New()
	viperInstance.SetConfigName("config_test.yaml")
	viperInstance.SetConfigType("yaml")
	viperInstance.AddConfigPath(".")
	err := viperInstance.ReadInConfig()
	assert.NoError(t, err)

	c1, err := LoadConfig(*viperInstance, "./config_test.yaml")
	assert.NoError(t, err)

	c2 := Config{
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
