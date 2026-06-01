package report

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type Input struct {
	ObjectName    string
	ObjectType    string
	Summary       string
	MarkdownBody  string
	Mermaid       string
	CreatedAt     time.Time
	Risk          string
	TokensUsed    int
}

func Markdown(in Input) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", html.EscapeString(in.ObjectName))
	fmt.Fprintf(&sb, "**Type:** %s  \n", in.ObjectType)
	fmt.Fprintf(&sb, "**Risk:** %s  \n", in.Risk)
	fmt.Fprintf(&sb, "**Generated:** %s  \n", in.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "**Tokens:** %d\n\n", in.TokensUsed)
	fmt.Fprintf(&sb, "## Summary\n\n%s\n\n", html.EscapeString(in.Summary))
	fmt.Fprintf(&sb, "## Diagram\n\n```mermaid\n%s\n```\n\n", in.Mermaid)
	if in.MarkdownBody != "" {
		sb.WriteString(in.MarkdownBody)
		sb.WriteString("\n")
	}
	return sb.String()
}

func HTML(in Input) string {
	body := Markdown(in)
	var sb strings.Builder
	fmt.Fprintf(&sb, `<!doctype html><html><head><meta charset="utf-8">`)
	fmt.Fprintf(&sb, `<title>%s</title>`, html.EscapeString(in.ObjectName))
	sb.WriteString(`<style>body{font-family:system-ui,sans-serif;max-width:860px;margin:2rem auto;padding:0 1rem;color:#1f2937;background:#fff}pre{background:#f3f4f6;padding:1rem;border-radius:6px;overflow:auto}code{font-family:ui-monospace,Menlo,monospace;font-size:0.9em}h1,h2,h3{color:#111827}a{color:#2563eb}</style></head><body>`)
	for _, line := range strings.Split(body, "\n") {
		switch {
		case strings.HasPrefix(line, "# "):
			fmt.Fprintf(&sb, "<h1>%s</h1>", html.EscapeString(strings.TrimPrefix(line, "# ")))
		case strings.HasPrefix(line, "## "):
			fmt.Fprintf(&sb, "<h2>%s</h2>", html.EscapeString(strings.TrimPrefix(line, "## ")))
		case strings.HasPrefix(line, "### "):
			fmt.Fprintf(&sb, "<h3>%s</h3>", html.EscapeString(strings.TrimPrefix(line, "### ")))
		case strings.HasPrefix(line, "- "):
			fmt.Fprintf(&sb, "<li>%s</li>", html.EscapeString(strings.TrimPrefix(line, "- ")))
		case strings.HasPrefix(line, "```mermaid"):
			sb.WriteString("<pre class=\"mermaid\">")
		case line == "```":
			sb.WriteString("</pre>")
		case strings.HasPrefix(line, "```"):
			fmt.Fprintf(&sb, "<pre><code>%s</code></pre>", html.EscapeString(strings.TrimPrefix(line, "```")))
		default:
			if strings.TrimSpace(line) == "" {
				sb.WriteString("<br/>")
			} else {
				fmt.Fprintf(&sb, "<p>%s</p>", html.EscapeString(line))
			}
		}
	}
	sb.WriteString("</body></html>")
	return sb.String()
}
