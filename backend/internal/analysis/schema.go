package analysis

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidJSON = errors.New("invalid json")

var mermaidKeyword = regexp.MustCompile(`(?i)^(graph|flowchart)\s+(TD|LR|TB|RL|BT)?`)

type Schema struct {
	ObjectName                string   `json:"objectName"`
	ObjectType                string   `json:"objectType"`
	Summary                   string   `json:"summary"`
	SimpleExplanation         string   `json:"simpleExplanation"`
	ExecutionFlow             []string `json:"executionFlow"`
	MermaidDiagram            string   `json:"mermaidDiagram"`
	Tables                    []any    `json:"tables"`
	Parameters                []any    `json:"parameters"`
	BusinessRules             []any    `json:"businessRules"`
	Risks                     []any    `json:"risks"`
	PossibleBugs              []any    `json:"possibleBugs"`
	ModernizationSuggestions  []any    `json:"modernizationSuggestions"`
	MarkdownReport            string   `json:"markdownReport"`
}

func ValidateAndParse(raw string) (*Schema, error) {
	cleaned := stripCodeFence(raw)
	var s Schema
	if err := json.Unmarshal([]byte(cleaned), &s); err != nil {
		return nil, ErrInvalidJSON
	}
	if s.ObjectName == "" {
		s.ObjectName = "UNKNOWN"
	}
	if s.ObjectType == "" {
		s.ObjectType = "unknown"
	}
	if s.MermaidDiagram != "" && !looksLikeMermaid(s.MermaidDiagram) {
		return nil, errors.New("invalid mermaid diagram")
	}
	return &s, nil
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i > 0 {
			s = s[i+1:]
		}
		if strings.HasSuffix(s, "```") {
			s = s[:len(s)-3]
		}
	}
	return strings.TrimSpace(s)
}

func looksLikeMermaid(s string) bool {
	first := strings.TrimSpace(s)
	if first == "" {
		return false
	}
	return mermaidKeyword.MatchString(first)
}

func RepairJSON(raw string) string {
	cleaned := stripCodeFence(raw)
	cleaned = strings.TrimSpace(cleaned)
	if !strings.HasPrefix(cleaned, "{") {
		start := strings.Index(cleaned, "{")
		if start >= 0 {
			cleaned = cleaned[start:]
		}
	}
	if !strings.HasSuffix(cleaned, "}") {
		end := strings.LastIndex(cleaned, "}")
		if end >= 0 {
			cleaned = cleaned[:end+1]
		}
	}
	return cleaned
}
