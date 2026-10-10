package registrydocs

import (
	"strings"
	"testing"
)

func TestMarkdown_searchEngineRequiresQuery(t *testing.T) {
	text := Markdown("search_engine")
	if !strings.Contains(text, "query") || !strings.Contains(text, "scraped or search text") || strings.Contains(text, "UNTRUSTED_WEB_CONTENT") {
		t.Fatal("search guide missing")
	}
	if Markdown("no_such_tool") != "" {
		t.Fatal("missing tool returned text")
	}
}
