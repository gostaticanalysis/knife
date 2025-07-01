package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCutterHandlerWithComplexTemplate(t *testing.T) {
	cases := []struct {
		name         string
		format       string
		wantErr      bool
		wantContains string
	}{
		{
			name:    "simple template",
			format:  "{{.Name}}",
			wantErr: false,
		},
		{
			name:    "range template",
			format:  `{{range .Funcs}}{{.Name}}{{end}}`,
			wantErr: false,
		},
		{
			name:    "complex template with if",
			format:  `{{range .Funcs}}{{if eq .Name "Execute"}}found{{end}}{{end}}`,
			wantErr: false,
		},
		{
			name:    "template with no match - empty result",
			format:  `{{range .Funcs}}{{if eq .Name "NonExistentFunction"}}found{{end}}{{end}}`,
			wantErr: false,
		},
		{
			name:    "issue 34 template",
			format:  `{{range .Funcs}}{{if eq .Name "Execute"}}{{.Name}}: {{range .Signature.Params}}{{.Name}} {{.Type}}{{end}}{{end}}{{end}}`,
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := CutterInput{
				Patterns: []string{"."},
				Format:   tc.format,
			}

			// Create the params
			params := &mcpsdk.CallToolParamsFor[CutterInput]{
				Name:      "cutter",
				Arguments: input,
			}

			// Call the handler directly
			result, err := cutterHandler(context.Background(), nil, params)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if result == nil {
					t.Error("expected result but got nil")
				} else {
					// Check if result contains valid content
					if len(result.Content) == 0 {
						t.Error("expected content in result but got empty")
					} else {
						// Parse JSON response
						textContent, ok := result.Content[0].(*mcpsdk.TextContent)
						if !ok {
							t.Errorf("expected TextContent but got %T", result.Content[0])
						} else {
							var output CutterOutput
							if err := json.Unmarshal([]byte(textContent.Text), &output); err != nil {
								t.Errorf("failed to parse JSON response: %v", err)
							} else {
								if !output.Success {
									t.Errorf("expected success=true but got success=false, error: %s", output.Error)
								}
								if len(output.Results) == 0 {
									t.Error("expected results but got empty")
								}
								t.Logf("Success: %t, Results count: %d", output.Success, len(output.Results))
							}
						}
					}
				}
			}
		})
	}
}
