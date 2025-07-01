package mcp

import (
	"bytes"
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gostaticanalysis/knife"
)

// KnifeInput represents the input parameters for the knife MCP tool.
// It defines the package patterns to analyze, optional template formatting,
// extra data for template evaluation, and XPath filtering for AST nodes.
type KnifeInput struct {
	Patterns []string `json:"patterns"`         // Package patterns to analyze (e.g., ["fmt", "net/http"])
	Format   string   `json:"format,omitempty"` // Template string for output formatting
	Data     string   `json:"data,omitempty"`   // Extra data as key:value pairs
	XPath    string   `json:"xpath,omitempty"`  // XPath expression for AST node filtering
}

// KnifeOutput represents the output from the knife MCP tool.
// It contains the formatted analysis results as a string.
type KnifeOutput struct {
	Result string `json:"result"` // Formatted analysis results
}

// newKnifeTool creates the knife MCP tool.
func newKnifeTool() *mcp.ServerTool {
	description := "Execute knife to analyze Go packages with template formatting."
	formatDesc := fmt.Sprintf("Template string for output formatting (default: \"{{.}}\").\n\n%s", templateDoc)

	return mcp.NewServerTool("knife", description, knifeHandler,
		mcp.Input(
			mcp.Property("patterns", mcp.Description("Package patterns to analyze (e.g., [\"fmt\", \"net/http\", \"./...\"])"), mcp.Required(true)),
			mcp.Property("format", mcp.Description(formatDesc)),
			mcp.Property("data", mcp.Description("Extra data as key:value pairs (e.g., \"key1:value1,key2:value2\")")),
			mcp.Property("xpath", mcp.Description("XPath expression for AST node filtering")),
		),
	)
}

// knifeHandler handles the knife tool execution.
func knifeHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[KnifeInput]) (*mcp.CallToolResultFor[KnifeOutput], error) {
	input := params.Arguments

	if len(input.Patterns) == 0 {
		return nil, fmt.Errorf("patterns is required")
	}

	// Create knife instance
	knifeOpt := &knife.KnifeOption{Tests: true}
	k, err := knife.New(knifeOpt, input.Patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to create knife: %w", err)
	}

	// Set up options
	opt := &knife.ExecuteOption{
		XPath: input.XPath,
	}

	// Parse extra data if provided
	if input.Data != "" {
		extraData, err := parseExtraData(input.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse data: %w", err)
		}
		opt.ExtraData = extraData
	}

	// Use default format if not provided
	format := input.Format
	if format == "" {
		format = "{{.}}"
	}

	// Execute knife for each package
	var buf bytes.Buffer
	pkgs := k.Packages()
	for i, pkg := range pkgs {
		if err := k.Execute(&buf, pkg, format, opt); err != nil {
			return nil, fmt.Errorf("failed to execute knife for package %s: %w", pkg.PkgPath, err)
		}

		if i != len(pkgs)-1 {
			fmt.Fprintln(&buf)
		}
	}

	return &mcp.CallToolResultFor[KnifeOutput]{
		Content: []mcp.Content{&mcp.TextContent{
			Text: buf.String(),
		}},
	}, nil
}
