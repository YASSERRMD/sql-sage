package report

import (
	"strings"
	"testing"
	"time"
)

func sample() Input {
	return Input{
		ObjectName:   "close_account",
		ObjectType:   "procedure",
		Summary:      "Closes an account if balance is zero.",
		MarkdownBody: "## Steps\n- read\n- update",
		Mermaid:      "flowchart TD\nA-->B",
		CreatedAt:    time.Now(),
		Risk:         "low",
		TokensUsed:   123,
	}
}

func TestMarkdown(t *testing.T) {
	out := Markdown(sample())
	if !strings.Contains(out, "# close_account") {
		t.Fatal("missing title")
	}
	if !strings.Contains(out, "```mermaid") {
		t.Fatal("missing mermaid block")
	}
}

func TestHTML(t *testing.T) {
	out := HTML(sample())
	if !strings.Contains(out, "<h1>close_account</h1>") {
		t.Fatal("missing h1")
	}
	if !strings.Contains(out, "<!doctype html>") {
		t.Fatal("missing doctype")
	}
}
