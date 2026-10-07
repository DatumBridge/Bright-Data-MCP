package registrydocs

import (
	"embed"
	"strings"
)

//go:embed *.md
var files embed.FS

// Markdown returns the Tool Registry guide for one tool.
// Deploy and Setup copy this text from tools/list into the Docs field.
func Markdown(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return ""
	}
	raw, err := files.ReadFile(name + ".md")
	if err != nil {
		return ""
	}
	text := strings.TrimSpace(string(raw))
	const maxDocs = 24000
	if len(text) > maxDocs {
		text = text[:maxDocs]
	}
	return text
}
