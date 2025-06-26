package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewKnifeTool(t *testing.T) {
	tool := newKnifeTool()
	if tool == nil {
		t.Fatal("newKnifeTool returned nil")
	}

	if tool.Tool.Name != "knife" {
		t.Errorf("expected tool name 'knife', got %s", tool.Tool.Name)
	}
}

func TestParseExtraData(t *testing.T) {
	cases := []struct {
		input    string
		expected map[string]any
		hasError bool
	}{
		{
			input:    "key1:value1,key2:value2",
			expected: map[string]any{"key1": "value1", "key2": "value2"},
			hasError: false,
		},
		{
			input:    "single:value",
			expected: map[string]any{"single": "value"},
			hasError: false,
		},
		{
			input:    "",
			expected: map[string]any{},
			hasError: false,
		},
		{
			input:    "invalid_format",
			expected: nil,
			hasError: true,
		},
		{
			input:    "key1:value1,invalid",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range cases {
		result, err := parseExtraData(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("expected error for input %q", tc.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("unexpected error for input %q: %v", tc.input, err)
			continue
		}

		if len(result) != len(tc.expected) {
			t.Errorf("expected %d entries, got %d for input %q", len(tc.expected), len(result), tc.input)
			continue
		}

		for k, v := range tc.expected {
			if result[k] != v {
				t.Errorf("expected %q:%q, got %q:%q for input %q", k, v, k, result[k], tc.input)
			}
		}
	}
}

func TestKnifeHandler_ValidationErrors(t *testing.T) {
	cases := []struct {
		name  string
		input KnifeInput
	}{
		{
			name:  "empty patterns",
			input: KnifeInput{Patterns: []string{}},
		},
		{
			name:  "nil patterns",
			input: KnifeInput{Patterns: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := &mcp.CallToolParamsFor[KnifeInput]{
				Arguments: tc.input,
			}

			_, err := knifeHandler(context.Background(), nil, params)
			if err == nil {
				t.Error("expected validation error")
			}
		})
	}
}
