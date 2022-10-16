package strings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testMetricStringer struct{}

func (t *testMetricStringer) String() string {
	return `test-stringer`
}

func TestMetric(t *testing.T) {
	jsonBytes, _ := json.Marshal(
		map[string]any{
			`data`: 5.89,
		},
	)

	testCases := []struct {
		chunks   []any
		expected string
	}{
		{
			chunks: []any{
				`internal`,
				`module`,
				`getPaidAccount`,
				`timings`,
			},
			expected: `internal.module.getPaidAccount.timings`,
		},
		{
			chunks: []any{
				`internal`,
				testMetricStringer{}, // !!! CAUTION !!! produces `{}` instead of `test-stringer`
				&testMetricStringer{},
				`success`,
			},
			expected: `internal.{}.test-stringer.success`,
		},
		{
			chunks: []any{
				`internal`,
				jsonBytes,
				string(jsonBytes),
				`success`,
			},
			expected: `internal.[123 34 100 97 116 97 34 58 53 46 56 57 125].{"data":5.89}.success`,
		},
	}
	for _, testCase := range testCases {
		t.Run(
			`test_metric:`+testCase.expected, func(t *testing.T) {
				actual := Metric(testCase.chunks...)
				require.Equal(t, testCase.expected, actual)
			},
		)
	}
}
