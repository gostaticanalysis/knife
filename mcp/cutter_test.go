package mcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewCutterTool(t *testing.T) {
	tool := newCutterTool()
	if tool == nil {
		t.Fatal("newCutterTool returned nil")
	}

	if tool.Tool.Name != "cutter" {
		t.Errorf("expected tool name 'cutter', got %s", tool.Tool.Name)
	}
}

func TestCutterHandler_ValidationErrors(t *testing.T) {
	cases := []struct {
		name  string
		input CutterInput
	}{
		{
			name:  "empty patterns",
			input: CutterInput{Patterns: []string{}},
		},
		{
			name:  "nil patterns",
			input: CutterInput{Patterns: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := &mcp.CallToolParamsFor[CutterInput]{
				Arguments: tc.input,
			}

			_, err := cutterHandler(context.Background(), nil, params)
			if err == nil {
				t.Error("expected validation error")
			}
		})
	}
}
