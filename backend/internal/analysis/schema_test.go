package analysis

import (
	"strings"
	"testing"
)

func TestValidateAndParseOK(t *testing.T) {
	raw := `{"objectName":"X","objectType":"procedure","summary":"s","simpleExplanation":"e","executionFlow":["a"],"mermaidDiagram":"flowchart TD\nA-->B","tables":[],"parameters":[],"businessRules":[],"risks":[],"possibleBugs":[],"modernizationSuggestions":[],"markdownReport":"# X"}`
	s, err := ValidateAndParse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if s.ObjectName != "X" {
		t.Fatalf("expected X got %s", s.ObjectName)
	}
}

func TestValidateAndParseInvalid(t *testing.T) {
	if _, err := ValidateAndParse("not json"); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateAndParseWithFence(t *testing.T) {
	raw := "```json\n{\"objectName\":\"X\"}\n```"
	s, err := ValidateAndParse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if s.ObjectName != "X" {
		t.Fatalf("expected X got %s", s.ObjectName)
	}
}

func TestMermaidInvalid(t *testing.T) {
	raw := `{"objectName":"X","mermaidDiagram":"hello world"}`
	if _, err := ValidateAndParse(raw); err == nil {
		t.Fatal("expected mermaid error")
	}
}

func TestRepairJSON(t *testing.T) {
	got := RepairJSON("noise {\"a\":1} more noise")
	if !strings.HasPrefix(got, "{") || !strings.HasSuffix(got, "}") {
		t.Fatalf("expected trimmed, got %q", got)
	}
}
