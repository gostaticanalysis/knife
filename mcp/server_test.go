package mcp

import (
	"strings"
	"testing"
)

func TestNewKnifeServer(t *testing.T) {
	server := NewKnifeServer()
	if server == nil {
		t.Fatal("NewKnifeServer returned nil")
	}
}

func TestNewCutterServer(t *testing.T) {
	server := NewCutterServer()
	if server == nil {
		t.Fatal("NewCutterServer returned nil")
	}
}

func TestTemplateDocEmbedded(t *testing.T) {
	if templateDoc == "" {
		t.Fatal("templateDoc is empty, embed failed")
	}

	if len(templateDoc) < 100 {
		t.Fatalf("templateDoc seems too short: %d chars", len(templateDoc))
	}

	// Check if it contains expected content
	if !strings.Contains(templateDoc, "Template Format Reference") {
		t.Error("templateDoc does not contain expected title")
	}
}
