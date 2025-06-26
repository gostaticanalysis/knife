package mcp

import (
	"bytes"
	"context"
	"fmt"

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
// It contains the formatted analysis results as a string.
type CutterOutput struct {
	Result string `json:"result"` // Formatted analysis results
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
	c, err := cutter.New(input.Patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cutter: %w", err)
	}

	// Parse extra data if provided
	extraData, err := parseExtraData(input.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// Use default format if not provided
	format := input.Format
	if format == "" {
		format = "{{.}}"
	}

	// Execute cutter for each package
	var buf bytes.Buffer
	pkgs := c.KnifePackages()
	for i, pkg := range pkgs {
		opt := &cutter.Option{ExtraData: extraData}
		if err := c.Execute(&buf, pkg, format, opt); err != nil {
			return nil, fmt.Errorf("failed to execute cutter for package %s: %w", pkg.Path, err)
		}

		if i != len(pkgs)-1 {
			fmt.Fprintln(&buf)
		}
	}

	return &mcp.CallToolResultFor[CutterOutput]{
		Content: []mcp.Content{&mcp.TextContent{
			Text: buf.String(),
		}},
	}, nil
}
