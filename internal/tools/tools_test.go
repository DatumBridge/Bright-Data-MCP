package tools_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/mcp"
	"github.com/datumbridge/bright-data-mcp/internal/tools"
)

func clearToolFilterEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"BRIGHTDATA_PRO_MODE", "PRO_MODE",
		"BRIGHTDATA_GROUPS", "GROUPS",
		"BRIGHTDATA_TOOLS", "TOOLS",
	} {
		t.Setenv(k, "")
	}
}

func TestRegisterToolNames(t *testing.T) {
	clearToolFilterEnv(t)
	descs, handlers := tools.Register()
	want := map[string]bool{
		"search_engine":       false,
		"scrape_as_markdown":  false,
		"search_engine_batch": false,
		"scrape_batch":        false,
		"scrape_url":          false,
		"search_serp":         false,
		"bright_data_health":  false,
	}
	if len(descs) != 7 {
		t.Fatalf("expected 7 rapid-mode tools, got %d", len(descs))
	}
	for _, d := range descs {
		if d.Name == "discover" {
			t.Fatalf("deprecated discover must not be registered")
		}
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

func TestRegisterProModeIncludesWebData(t *testing.T) {
	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_PRO_MODE", "true")
	descs, _ := tools.Register()
	if len(descs) < 60 {
		t.Fatalf("expected 60+ tools in pro mode, got %d", len(descs))
	}
	for _, d := range descs {
		if d.Name == "discover" {
			t.Fatalf("deprecated discover must not be registered in pro mode")
		}
	}
	required := []string{
		"search_engine", "scrape_as_markdown", "search_engine_batch", "scrape_batch",
		"scrape_as_html", "extract", "session_stats",
		"list_dataset_fields", "search_dataset",
		"web_data_amazon_product", "web_data_linkedin_person_profile",
		"web_data_npm_package", "web_data_chatgpt_ai_insights",
		"web_data_booking_hotel_listings", "web_data_reuter_news", "web_data_reddit_comments",
		"scraping_browser_navigate", "scraping_browser_snapshot",
	}
	have := map[string]bool{}
	for _, d := range descs {
		have[d.Name] = true
	}
	for _, name := range required {
		if !have[name] {
			t.Fatalf("pro mode missing official tool %s", name)
		}
	}
}

func TestSessionStatsNotInRapid(t *testing.T) {
	clearToolFilterEnv(t)
	descs, _ := tools.Register()
	for _, d := range descs {
		if d.Name == "session_stats" {
			t.Fatalf("session_stats must be Pro/advanced_scraping only, not Rapid")
		}
	}
}

func toolNames(descs []mcp.ToolDesc) map[string]bool {
	have := map[string]bool{}
	for _, d := range descs {
		have[d.Name] = true
	}
	return have
}

func TestGroupsBusinessTravelResearchSocial(t *testing.T) {
	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "business")
	have := toolNames(mustDescs(t))
	if have["discover"] {
		t.Fatal("business group must not include discover")
	}
	for _, name := range []string{
		"search_engine", "scrape_as_markdown",
		"web_data_crunchbase_company", "web_data_zoominfo_company_profile",
		"web_data_google_maps_reviews", "web_data_zillow_properties_listing",
		"list_dataset_fields", "search_dataset",
	} {
		if !have[name] {
			t.Fatalf("business missing %s", name)
		}
	}
	if have["web_data_booking_hotel_listings"] {
		t.Fatal("business must not include booking (travel only)")
	}

	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "travel")
	have = toolNames(mustDescs(t))
	if !have["web_data_booking_hotel_listings"] {
		t.Fatal("travel must include booking")
	}

	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "research")
	have = toolNames(mustDescs(t))
	if !have["web_data_github_repository_file"] {
		t.Fatal("research missing github file tool")
	}
	if have["web_data_reuter_news"] {
		t.Fatal("research must not force reuter_news")
	}

	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "social")
	have = toolNames(mustDescs(t))
	if !have["web_data_reddit_posts"] {
		t.Fatal("social missing reddit_posts")
	}
	if have["web_data_reddit_comments"] {
		t.Fatal("social must not include reddit_comments (Pro extra only)")
	}
	if have["discover"] {
		t.Fatal("social must not include discover")
	}

	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "advanced_scraping")
	have = toolNames(mustDescs(t))
	for _, name := range []string{
		"search_engine_batch", "scrape_batch", "scrape_as_html", "extract", "session_stats",
	} {
		if !have[name] {
			t.Fatalf("advanced_scraping missing %s", name)
		}
	}
}

func TestGroupsEcommerceExcludesDiscover(t *testing.T) {
	clearToolFilterEnv(t)
	t.Setenv("BRIGHTDATA_GROUPS", "ecommerce")
	have := toolNames(mustDescs(t))
	if have["discover"] {
		t.Fatal("ecommerce must not register discover")
	}
	if !have["web_data_amazon_product"] || !have["search_engine"] {
		t.Fatal("ecommerce missing expected tools")
	}
}

func mustDescs(t *testing.T) []mcp.ToolDesc {
	t.Helper()
	descs, _ := tools.Register()
	return descs
}

func TestScrapeBatchEmptyAndCap(t *testing.T) {
	clearToolFilterEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	empty := handlers["scrape_batch"](mustJSON(map[string]interface{}{
		"urls": []interface{}{},
	}))
	if empty["isError"] != true {
		t.Fatalf("expected empty urls error: %#v", empty)
	}

	urls := make([]interface{}, 0, 12)
	for i := 0; i < 12; i++ {
		urls = append(urls, fmt.Sprintf("https://example.com/p%d", i))
	}
	res := handlers["scrape_batch"](mustJSON(map[string]interface{}{"urls": urls}))
	if res["isError"] != false {
		t.Fatalf("expected success with capped urls: %#v", res)
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
		"format": "raw",
	}))
	if res["isError"] != false {
		t.Fatalf("expected success: %#v", res)
	}
	text := contentText(res)
	if !strings.Contains(text, "UNTRUSTED_WEB_CONTENT") || !strings.Contains(text, "# Page") {
		t.Fatalf("unexpected text: %s", text)
	}
	if !strings.Contains(text, "--- content ---") {
		t.Fatalf("expected content section: %s", text)
	}
}

func TestScrapeURLDefaultFormatIsRaw(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<html></html>"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["scrape_url"](mustJSON(map[string]interface{}{
		"url": "https://example.com",
	}))
	if res["isError"] != false {
		t.Fatalf("expected success: %#v", res)
	}
	if got["format"] != "raw" {
		t.Fatalf("expected format=raw, got %#v", got)
	}
	if _, ok := got["data_format"]; ok {
		t.Fatalf("expected no data_format by default, got %#v", got)
	}
}

func TestScrapeURLMarkdownUsesDataFormat(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# Page"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["scrape_url"](mustJSON(map[string]interface{}{
		"url":         "https://example.com",
		"format":      "raw",
		"data_format": "markdown",
	}))
	if res["isError"] != false {
		t.Fatalf("expected success: %#v", res)
	}
	if got["format"] != "raw" || got["data_format"] != "markdown" {
		t.Fatalf("unexpected upstream payload: %#v", got)
	}
}

func TestSearchEngineReturnsMarkdownWithoutParsedJSON(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# LANDCO\n\n[Site](https://landco.example)"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["search_engine"](mustJSON(map[string]interface{}{
		"query": "LANDCO",
	}))
	if res["isError"] != false {
		t.Fatalf("markdown search must succeed: %#v", res)
	}
	if got["data_format"] != "markdown" {
		t.Fatalf("expected markdown, got %#v", got)
	}
	if url, _ := got["url"].(string); strings.Contains(url, "brd_json=1") {
		t.Fatalf("parsed JSON requested: %s", url)
	}
	text := contentText(res)
	if !strings.Contains(text, "# LANDCO") || !strings.Contains(text, "https://landco.example") {
		t.Fatalf("missing page text: %s", text)
	}
}

func TestSearchEngineKeepsOrganicJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"organic":[{"title":"LANDCO","link":"https://landco.example","description":"co"}]}`))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["search_engine"](mustJSON(map[string]interface{}{"query": "LANDCO"}))
	if res["isError"] != false {
		t.Fatalf("%#v", res)
	}
	text := contentText(res)
	if !strings.Contains(text, `"title": "LANDCO"`) {
		t.Fatalf("expected organic JSON: %s", text)
	}
}

func TestSearchSERPDefaultSkipsBrdJSON(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<html>serp</html>"))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_SERP_ZONE", "serp1")

	_, handlers := tools.Register()
	res := handlers["search_serp"](mustJSON(map[string]interface{}{"query": "LANDCO"}))
	if res["isError"] != false {
		t.Fatalf("%#v", res)
	}
	url, _ := got["url"].(string)
	if strings.Contains(url, "brd_json=1") {
		t.Fatalf("default request asked for parsed JSON: %s", url)
	}
	text := contentText(res)
	if !strings.Contains(text, "raw_excerpt") || !strings.Contains(text, "serp") {
		t.Fatalf("expected raw page in result: %s", text)
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
		"brd_json":    true,
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

func TestScrapeAsMarkdownSurfacesBrightDataResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Brd-Error", "no_peer")
		w.Header().Set("X-Brd-Error-Code", "client_10100")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	t.Setenv("BRIGHTDATA_API_KEY", "k")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "wu1")

	_, handlers := tools.Register()
	res := handlers["scrape_as_markdown"](mustJSON(map[string]interface{}{
		"url": "https://example.com/about",
	}))
	if res["isError"] != false {
		t.Fatalf("HTTP 200 must still return the response text: %#v", res)
	}
	text := contentText(res)
	for _, want := range []string{
		`"zone": "wu1"`,
		`"data_format": "markdown"`,
		`"empty_body": true`,
		`"char_count": 0`,
		`"success": false`,
		`"x-brd-error": "no_peer"`,
		`"x-brd-error-code": "client_10100"`,
		"--- content ---",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
	if !strings.HasSuffix(text, "--- content ---\n\n") {
		t.Fatalf("body was rewritten: %q", text)
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
