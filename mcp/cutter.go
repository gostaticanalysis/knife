package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gostaticanalysis/knife/cutter"
)

// CutterInput represents the input parameters for the cutter MCP tool.
// It defines the package patterns to analyze, optional template formatting,
// and extra data for template evaluation.
type CutterInput struct {
	Patterns []string `json:"patterns"`         // Package patterns to analyze (e.g., ["fmt", "net/http"])
	Format   string   `json:"format,omitempty"` // Template string for output formatting
	Data     string   `json:"data,omitempty"`   // Extra data as key:value pairs
}

// CutterOutput represents the output from the cutter MCP tool.
// It contains the formatted analysis results as structured JSON.
type CutterOutput struct {
	Success bool            `json:"success"`         // Whether the operation succeeded
	Results []PackageResult `json:"results"`         // Analysis results per package
	Error   string          `json:"error,omitempty"` // Error message if any
}

// PackageResult represents the analysis result for a single package.
type PackageResult struct {
	PackageName string `json:"package_name"` // Name of the analyzed package
	Content     string `json:"content"`      // Formatted template output
}

// newCutterTool creates the cutter MCP tool.
func newCutterTool() *mcp.ServerTool {
	description := "Execute cutter to analyze Go package types with template formatting."
	formatDesc := fmt.Sprintf("Template string for output formatting (default: \"{{.}}\").\n\n%s", templateDoc)

	return mcp.NewServerTool("cutter", description, cutterHandler,
		mcp.Input(
			mcp.Property("patterns", mcp.Description("Package patterns to analyze (e.g., [\"fmt\", \"net/http\", \"./...\"])"), mcp.Required(true)),
			mcp.Property("format", mcp.Description(formatDesc)),
			mcp.Property("data", mcp.Description("Extra data as key:value pairs (e.g., \"key1:value1,key2:value2\")")),
		),
	)
}

// cutterHandler handles the cutter tool execution.
func cutterHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[CutterInput]) (*mcp.CallToolResultFor[CutterOutput], error) {
	input := params.Arguments

	if len(input.Patterns) == 0 {
		return nil, fmt.Errorf("patterns is required")
	}

	// Create cutter instance
	cutterOpt := &cutter.CutterOption{Tests: true}
	c, err := cutter.New(cutterOpt, input.Patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cutter: %w", err)
	}

	// Parse extra data if provided
	extraData, err := parseExtraData(input.Data)
	if err != nil {
		return &mcp.CallToolResultFor[CutterOutput]{
			Content: []mcp.Content{&mcp.TextContent{
				Text: mustMarshalJSON(CutterOutput{
					Success: false,
					Error:   fmt.Sprintf("failed to parse data: %q", err.Error()),
				}),
			}},
		}, nil
	}

	// Use default format if not provided
	format := input.Format
	if format == "" {
		format = "{{.}}"
	}

	// Execute cutter for each package
	pkgs := c.KnifePackages()
	results := make([]PackageResult, 0, len(pkgs))

	for _, pkg := range pkgs {
		var buf bytes.Buffer
		opt := &cutter.Option{ExtraData: extraData}
		if err := c.Execute(&buf, pkg, format, opt); err != nil {
			return &mcp.CallToolResultFor[CutterOutput]{
				Content: []mcp.Content{&mcp.TextContent{
					Text: mustMarshalJSON(CutterOutput{
						Success: false,
						Error:   fmt.Sprintf("failed to execute cutter for package %s: %q", pkg.Path, err.Error()),
					}),
				}},
			}, nil
		}

		content := buf.String()

		results = append(results, PackageResult{
			PackageName: pkg.Path,
			Content:     content,
		})
	}

	output := CutterOutput{
		Success: true,
		Results: results,
	}

	return &mcp.CallToolResultFor[CutterOutput]{
		Content: []mcp.Content{&mcp.TextContent{
			Text: mustMarshalJSON(output),
		}},
	}, nil
}

// mustMarshalJSON marshals v to JSON, panicking on error (should never happen with our types)
func mustMarshalJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// This should never happen with our defined types
		return fmt.Sprintf(`{"success": false, "error": "JSON marshal error: %q"}`, err.Error())
	}
	return string(b)
}
