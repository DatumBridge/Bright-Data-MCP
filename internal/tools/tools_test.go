package tools_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/tools"
)

func TestRegisterToolNames(t *testing.T) {
	descs, handlers := tools.Register()
	want := map[string]bool{
		"scrape_url":         false,
		"search_serp":        false,
		"bright_data_health": false,
	}
	if len(descs) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(descs))
	}
	for _, d := range descs {
		if _, ok := want[d.Name]; !ok {
			t.Fatalf("unexpected tool %s", d.Name)
		}
		want[d.Name] = true
		if _, ok := handlers[d.Name]; !ok {
			t.Fatalf("missing handler %s", d.Name)
		}
	}
	for name, seen := range want {
		if !seen {
			t.Fatalf("missing tool %s", name)
		}
	}
}

func TestScrapeURLMissingZone(t *testing.T) {
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "")
	t.Setenv("BRIGHTDATA_SERP_ZONE", "")
	_, handlers := tools.Register()
	res := handlers["scrape_url"](mustJSON(map[string]interface{}{
		"url": "https://example.com",
	}))
	if res["isError"] != true {
		t.Fatalf("expected error: %#v", res)
	}
}

func TestScrapeURLSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# Page"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["scrape_url"](mustJSON(map[string]interface{}{
		"url":    "https://example.com",
		"format": "markdown",
	}))
	if res["isError"] != false {
		t.Fatalf("expected success: %#v", res)
	}
	text := contentText(res)
	if !strings.Contains(text, "UNTRUSTED_WEB_CONTENT") || !strings.Contains(text, "# Page") {
		t.Fatalf("unexpected text: %s", text)
	}
}

func TestSearchSERPParsesOrganic(t *testing.T) {
	payload := `{"organic":[{"title":"A","link":"https://a.example","snippet":"sa"},{"title":"B","link":"https://b.example","snippet":"sb"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_SERP_ZONE", "serp1")

	_, handlers := tools.Register()
	res := handlers["search_serp"](mustJSON(map[string]interface{}{
		"query":       "test",
		"max_results": 1,
	}))
	if res["isError"] != false {
		t.Fatalf("expected success: %#v", res)
	}
	text := contentText(res)
	if !strings.Contains(text, `"title": "A"`) || strings.Contains(text, `"title": "B"`) {
		t.Fatalf("max_results not applied: %s", text)
	}
}

func TestHealthMasksKey(t *testing.T) {
	t.Setenv("BRIGHTDATA_API_KEY", "abcdefghijklmnop")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu")
	t.Setenv("BRIGHTDATA_SERP_ZONE", "serp")
	_, handlers := tools.Register()
	res := handlers["bright_data_health"](mustJSON(map[string]interface{}{}))
	if res["isError"] != false {
		t.Fatalf("%#v", res)
	}
	text := contentText(res)
	if strings.Contains(text, "abcdefghijklmnop") {
		t.Fatal("api key leaked")
	}
	if !strings.Contains(text, "api_key_masked") {
		t.Fatalf("%s", text)
	}
}

func mustJSON(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func contentText(res map[string]interface{}) string {
	content, _ := res["content"].([]map[string]string)
	if len(content) == 0 {
		// JSON unmarshaling of tool result uses []interface{}
		raw, _ := res["content"].([]interface{})
		if len(raw) == 0 {
			return ""
		}
		m, _ := raw[0].(map[string]interface{})
		s, _ := m["text"].(string)
		return s
	}
	return content[0]["text"]
}
