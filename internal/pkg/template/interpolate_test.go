package template

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpolate(t *testing.T) {
	testCases := []struct {
		name         string
		source       string
		replacements map[string]any
		expected     string
		wantError    bool
	}{
		{
			name:   `simple`,
			source: `Hello, {{ .user }}!`,
			replacements: map[string]any{
				"user": "John",
			},
			expected: `Hello, John!`,
		},
		{
			name:   `urlquery`,
			source: `Follow the link: {{ .link | urlquery }}`,
			replacements: map[string]any{
				"link": "https://ru.wikipedia.org/wiki/Интерполяция",
			},
			expected: `Follow the link: https%3A%2F%2Fru.wikipedia.org%2Fwiki%2F%D0%98%D0%BD%D1%82%D0%B5%D1%80%D0%BF%D0%BE%D0%BB%D1%8F%D1%86%D0%B8%D1%8F`,
		},
		{
			name:   `missingkey`,
			source: `key is {{ .missing }}`,
			replacements: map[string]any{
				"name": "John",
			},
			wantError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				actual, err := Interpolate(testCase.source, testCase.replacements)
				if testCase.wantError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, testCase.expected, actual)
				}
			},
		)
	}
}
