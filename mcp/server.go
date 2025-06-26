// Package mcp provides Model Context Protocol (MCP) server implementations
// for knife and cutter tools, enabling remote execution of Go static analysis
// operations through MCP-compatible clients.
package mcp

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gostaticanalysis/knife"
)

//go:embed template.md
var templateDoc string

// NewKnifeServer creates a new MCP server for the knife tool.
// The server provides remote access to knife's Go package analysis capabilities
// through the Model Context Protocol, allowing clients to execute knife commands
// with template formatting and XPath filtering.
func NewKnifeServer() *mcp.Server {
	server := mcp.NewServer("knife", knife.Version(), nil)
	server.AddTools(newKnifeTool())
	return server
}

// NewCutterServer creates a new MCP server for the cutter tool.
// The server provides remote access to cutter's Go package type analysis
// capabilities through the Model Context Protocol, allowing clients to
// execute cutter commands with template formatting.
func NewCutterServer() *mcp.Server {
	server := mcp.NewServer("cutter", knife.Version(), nil)
	server.AddTools(newCutterTool())
	return server
}

// parseExtraData parses the extra data string into a map.
// The extraData string should be in the format "key1:value1,key2:value2".
// Returns an empty map if extraData is empty, or an error if the format is invalid.
func parseExtraData(extraData string) (map[string]any, error) {
	if extraData == "" {
		return make(map[string]any), nil
	}

	m := make(map[string]any)
	kvs := strings.Split(extraData, ",")
	for i := range kvs {
		key, value, ok := strings.Cut(strings.TrimSpace(kvs[i]), ":")
		if !ok {
			return nil, fmt.Errorf("invalid extra data format at '%s': expected 'key:value' pairs separated by commas", kvs[i])
		}
		m[key] = value
	}
	return m, nil
}
