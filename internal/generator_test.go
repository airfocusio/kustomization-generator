package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCombinedKubernetesResources1(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output []KustomizationResource
	}{
		{
			name:   "empty-1",
			input:  "",
			output: []KustomizationResource{},
		},
		{
			name:   "empty-2",
			input:  "\n",
			output: []KustomizationResource{},
		},
		{
			name:   "empty-3",
			input:  "---",
			output: []KustomizationResource{},
		},
		{
			name:   "empty-4",
			input:  "---\n",
			output: []KustomizationResource{},
		},
		{
			name:   "empty-5",
			input:  "\n---\n",
			output: []KustomizationResource{},
		},
		{
			name:  "single",
			input: mockResource("Secret", "database"),
			output: []KustomizationResource{
				{
					Name:    "database-secret.yaml",
					Content: mockResource("Secret", "database"),
				},
			},
		},
		{
			name:  "multiple",
			input: mockResource("Secret", "database") + "---\n" + mockResource("Secret", "other"),
			output: []KustomizationResource{
				{
					Name:    "database-secret.yaml",
					Content: mockResource("Secret", "database"),
				},
				{
					Name:    "other-secret.yaml",
					Content: mockResource("Secret", "other"),
				},
			},
		},
		{
			name:  "name-collision",
			input: mockResource("Secret", "database") + "---\n" + mockResource("Secret", "database"),
			output: []KustomizationResource{
				{
					Name:    "database-secret.yaml",
					Content: mockResource("Secret", "database"),
				},
				{
					Name:    "database-secret-1.yaml",
					Content: mockResource("Secret", "database"),
				},
			},
		},
		{
			name:   "comment-isolated",
			input:  "# foobar\n\n---\n",
			output: []KustomizationResource{},
		},
		{
			name:  "comment-before",
			input: "# foobar\n\n---\n" + mockResource("Secret", "database"),
			output: []KustomizationResource{
				{
					Name:    "database-secret.yaml",
					Content: "# foobar\n\n" + mockResource("Secret", "database"),
				},
			},
		},
		{
			name:  "comment-after",
			input: mockResource("Secret", "database") + "---\n# foobar\n\n",
			output: []KustomizationResource{
				{
					Name:    "database-secret.yaml",
					Content: mockResource("Secret", "database"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := splitCombinedKubernetesResources(testCase.input)
		if assert.NoError(t, err, "Case %s", testCase.name) {
			assert.Equal(t, testCase.output, actual, "Case %s", testCase.name)
		}
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

func mockResource(kind string, name string) string {
	return fmt.Sprintf(`apiVersion: v1
kind: %s
metadata:
  name: %s
`, kind, name)
}
