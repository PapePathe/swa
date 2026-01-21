package astformat

import (
	"encoding/json"
	"swahili/lang/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson_VisitFloatExpression(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		node         *ast.FloatExpression
		expected     map[string]any
		expectedJson string
		wantErr      bool
	}{
		{
			name:         "with 1.009",
			node:         &ast.FloatExpression{Value: float64(1.009)},
			expected:     map[string]any{"FloatExpression": 1.009},
			expectedJson: `{"FloatExpression":1.009}`,
		},
		{
			name:         "with 1.0",
			node:         &ast.FloatExpression{Value: float64(1.0)},
			expected:     map[string]any{"FloatExpression": 1.0},
			expectedJson: `{"FloatExpression":1}`,
		},
		{
			name:         "with 0.0",
			node:         &ast.FloatExpression{Value: float64(0.0)},
			expected:     map[string]any{"FloatExpression": 0.0},
			expectedJson: `{"FloatExpression":0}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var j Json = Json{Element: make(map[string]any)}

			gotErr := j.VisitFloatExpression(tt.node)

			assert.NoError(t, gotErr)

			actualJson, _ := json.Marshal(j.Element)

			assert.Equal(t, tt.expected, j.Element)
			assert.Equal(t, tt.expectedJson, string(actualJson))
		})
	}
}
