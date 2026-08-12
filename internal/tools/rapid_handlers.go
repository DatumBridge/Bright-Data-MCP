package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
	"github.com/datumbridge/bright-data-mcp/internal/mcp"
)

type toolAdder func(name, desc string, props map[string]interface{}, required []string, h mcp.ToolHandler)

func registerRapidTools(cfg ServerConfig, add toolAdder) {
	reg := func(name string, enabled bool, desc string, props map[string]interface{}, required []string, h mcp.ToolHandler) {
		if !enabled || !cfg.toolEnabled(name) {
			return
		}
		wrapped := func(raw json.RawMessage) map[string]interface{} {
			recordToolCall(name)
			return h(raw)
		}
		add(name, desc, props, required, wrapped)
	}

	reg("search_engine", true,
		"Scrape search results from Google, Bing or Yandex. Google returns JSON organic results; Bing/Yandex return Markdown.",
		baseProps(map[string]interface{}{
			"query":        map[string]interface{}{"type": "string"},
			"engine":       map[string]interface{}{"type": "string", "enum": []string{"google", "bing", "yandex"}, "default": "google"},
			"cursor":       map[string]interface{}{"type": "string", "description": "Pagination cursor (page index)"},
			"geo_location": map[string]interface{}{"type": "string", "description": "2-letter country code"},
		}),
		[]string{"query"}, handleSearchEngine)

	reg("scrape_as_markdown", true,
		"Scrape a single webpage via Web Unlocker and return Markdown.",
		baseProps(map[string]interface{}{
			"url": map[string]interface{}{"type": "string", "description": "Target URL"},
		}),
		[]string{"url"}, handleScrapeAsMarkdown)

	reg("scrape_as_html", true,
		"Scrape a single webpage via Web Unlocker and return raw HTML.",
		baseProps(map[string]interface{}{
			"url": map[string]interface{}{"type": "string"},
		}),
		[]string{"url"}, handleScrapeAsHTML)

	reg("search_engine_batch", true,
		"Run up to 10 search queries in parallel (Google JSON, Bing/Yandex Markdown).",
		baseProps(map[string]interface{}{
			"queries": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query":        map[string]interface{}{"type": "string"},
						"engine":       map[string]interface{}{"type": "string"},
						"cursor":       map[string]interface{}{"type": "string"},
						"geo_location": map[string]interface{}{"type": "string"},
					},
					"required": []string{"query"},
				},
				"maxItems": 10,
			},
		}),
		[]string{"queries"}, handleSearchEngineBatch)

	reg("scrape_batch", true,
		"Scrape up to 10 webpages and return URL/markdown pairs.",
		baseProps(map[string]interface{}{
			"urls": map[string]interface{}{
				"type": "array", "items": map[string]interface{}{"type": "string"}, "maxItems": 10,
			},
		}),
		[]string{"urls"}, handleScrapeBatch)

	reg("extract", true,
		"Scrape page as Markdown. Returns markdown plus optional extraction_prompt for your agent to structure as JSON (official MCP uses LLM sampling).",
		baseProps(map[string]interface{}{
			"url":               map[string]interface{}{"type": "string"},
			"extraction_prompt": map[string]interface{}{"type": "string"},
		}),
		[]string{"url"}, handleExtract)

	reg("discover", true,
		"AI-ranked web search via Bright Data Discover API.",
		baseProps(map[string]interface{}{
			"query":              map[string]interface{}{"type": "string"},
			"intent":             map[string]interface{}{"type": "string"},
			"country":            map[string]interface{}{"type": "string"},
			"city":               map[string]interface{}{"type": "string"},
			"language":           map[string]interface{}{"type": "string"},
			"num_results":        map[string]interface{}{"type": "integer"},
			"filter_keywords":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			"remove_duplicates":  map[string]interface{}{"type": "boolean"},
			"start_date":         map[string]interface{}{"type": "string"},
			"end_date":           map[string]interface{}{"type": "string"},
		}),
		[]string{"query"}, handleDiscover)

	reg("session_stats", true,
		"Report per-tool call counts for this MCP server process session.",
		baseProps(nil), nil, func(json.RawMessage) map[string]interface{} {
			return mcp.ToolResultText(sessionStatsText())
		})

	// Legacy DatumBridge tool names
	reg("scrape_url", true,
		"[Legacy alias] Use scrape_as_html or scrape_as_markdown. Bright Data Web Unlocker scrape.",
		baseProps(map[string]interface{}{
			"url":         map[string]interface{}{"type": "string"},
			"format":      map[string]interface{}{"type": "string", "enum": []string{"raw", "json"}, "default": "raw"},
			"data_format": map[string]interface{}{"type": "string", "enum": []string{"markdown", "screenshot"}},
			"country":     map[string]interface{}{"type": "string"},
			"max_chars":   map[string]interface{}{"type": "integer"},
		}),
		[]string{"url"}, handleScrapeURL)

	reg("search_serp", true,
		"[Legacy alias] Use search_engine. SERP via Web Unlocker.",
		baseProps(map[string]interface{}{
			"query":       map[string]interface{}{"type": "string"},
			"engine":      map[string]interface{}{"type": "string", "enum": []string{"google", "bing"}},
			"country":     map[string]interface{}{"type": "string"},
			"language":    map[string]interface{}{"type": "string"},
			"start":       map[string]interface{}{"type": "integer"},
			"brd_json":    map[string]interface{}{"type": "boolean"},
			"max_results": map[string]interface{}{"type": "integer"},
		}),
		[]string{"query"}, handleSearchSERP)

	reg("bright_data_health", true,
		"Validate Bright Data credentials/zones. probe=true runs a billable Unlocker request.",
		baseProps(map[string]interface{}{
			"probe": map[string]interface{}{"type": "boolean", "default": false},
		}), nil, handleHealth)
}

func unlockerZone(client *brightdata.Client) (string, error) {
	z := client.UnlockerZone()
	if z == "" {
		return "", fmt.Errorf("unlocker_zone required (credentials_json.unlocker_zone or BRIGHTDATA_UNLOCKER_ZONE)")
	}
	return z, nil
}

func handleSearchEngine(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone, err := unlockerZone(client)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	query := strArg(m, "query")
	engine := strArg(m, "engine")
	if engine == "" {
		engine = "google"
	}
	searchURL, err := brightdata.BuildEngineSearchURL(engine, query, strArg(m, "cursor"), strArg(m, "geo_location"))
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	isGoogle := engine == "google"
	opts := brightdata.RequestOpts{Zone: zone, URL: searchURL, Format: "raw"}
	if isGoogle {
		opts.URL = searchURL + "&brd_json=1"
		opts.DataFormat = "parsed_light"
	} else {
		opts.DataFormat = "markdown"
	}
	cctx, cancel := pollCtx()
	defer cancel()
	body, _, err := client.Request(cctx, opts)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if !isGoogle {
		return untrustedTextResult(string(body))
	}
	parsed, err := brightdata.ParseGoogleSearchResponse(body)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return jsonResult(parsed)
}

func handleScrapeAsMarkdown(raw json.RawMessage) map[string]interface{} {
	return scrapeUnlocker(raw, "markdown")
}

func handleScrapeAsHTML(raw json.RawMessage) map[string]interface{} {
	return scrapeUnlocker(raw, "")
}

func scrapeUnlocker(raw json.RawMessage, dataFormat string) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone, err := unlockerZone(client)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	target := strArg(m, "url")
	if err := brightdata.ValidateHTTPURL(target); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	opts := brightdata.RequestOpts{Zone: zone, URL: target, Format: "raw", DataFormat: dataFormat}
	cctx, cancel := pollCtx()
	defer cancel()
	body, status, err := client.Request(cctx, opts)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	content, truncated := brightdata.TruncateUTF8(string(body), brightdata.DefaultMaxChars())
	meta := map[string]interface{}{
		"success": true, "url": target, "status": status, "truncated": truncated,
		"char_count": len([]rune(content)),
	}
	return scrapeResult(meta, content)
}

func handleSearchEngineBatch(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone, err := unlockerZone(client)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	queries, ok := m["queries"].([]interface{})
	if !ok || len(queries) == 0 {
		return mcp.ToolResultError("queries array is required (max 10)")
	}
	if len(queries) > 10 {
		queries = queries[:10]
	}
	type result struct {
		Query  string      `json:"query"`
		Engine string      `json:"engine"`
		Result interface{} `json:"result,omitempty"`
		Error  string      `json:"error,omitempty"`
	}
	out := make([]result, 0, len(queries))
	for _, q := range queries {
		qm, _ := q.(map[string]interface{})
		query := strArg(qm, "query")
		engine := strArg(qm, "engine")
		if engine == "" {
			engine = "google"
		}
		searchURL, err := brightdata.BuildEngineSearchURL(engine, query, strArg(qm, "cursor"), strArg(qm, "geo_location"))
		if err != nil {
			out = append(out, result{Query: query, Engine: engine, Error: err.Error()})
			continue
		}
		opts := brightdata.RequestOpts{Zone: zone, URL: searchURL, Format: "raw"}
		if engine == "google" {
			opts.URL = searchURL + "&brd_json=1"
			opts.DataFormat = "parsed_light"
		} else {
			opts.DataFormat = "markdown"
		}
		cctx, cancel := pollCtx()
		body, _, reqErr := client.Request(cctx, opts)
		cancel()
		if reqErr != nil {
			out = append(out, result{Query: query, Engine: engine, Error: reqErr.Error()})
			continue
		}
		if engine == "google" {
			parsed, perr := brightdata.ParseGoogleSearchResponse(body)
			if perr != nil {
				out = append(out, result{Query: query, Engine: engine, Error: perr.Error()})
				continue
			}
			out = append(out, result{Query: query, Engine: engine, Result: parsed})
		} else {
			out = append(out, result{Query: query, Engine: engine, Result: string(body)})
		}
	}
	return jsonResult(out)
}

func handleScrapeBatch(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone, err := unlockerZone(client)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	urlsRaw, ok := m["urls"].([]interface{})
	if !ok || len(urlsRaw) == 0 {
		return mcp.ToolResultError("urls array is required (max 10)")
	}
	if len(urlsRaw) > 10 {
		urlsRaw = urlsRaw[:10]
	}
	type pair struct {
		URL     string `json:"url"`
		Content string `json:"content,omitempty"`
		Error   string `json:"error,omitempty"`
	}
	out := make([]pair, 0, len(urlsRaw))
	for _, u := range urlsRaw {
		target, _ := u.(string)
		target = strings.TrimSpace(target)
		opts := brightdata.RequestOpts{Zone: zone, URL: target, Format: "raw", DataFormat: "markdown"}
		cctx, cancel := pollCtx()
		body, _, reqErr := client.Request(cctx, opts)
		cancel()
		if reqErr != nil {
			out = append(out, pair{URL: target, Error: reqErr.Error()})
			continue
		}
		content, _ := brightdata.TruncateUTF8(string(body), brightdata.DefaultMaxChars())
		out = append(out, pair{URL: target, Content: content})
	}
	return jsonResult(out)
}

func handleExtract(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone, err := unlockerZone(client)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	target := strArg(m, "url")
	prompt := strArg(m, "extraction_prompt")
	cctx, cancel := pollCtx()
	defer cancel()
	body, _, err := client.Request(cctx, brightdata.RequestOpts{
		Zone: zone, URL: target, Format: "raw", DataFormat: "markdown",
	})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	meta := map[string]interface{}{
		"url": target, "note": "Structure the markdown below as JSON using your LLM per extraction_prompt",
	}
	if prompt != "" {
		meta["extraction_prompt"] = prompt
	}
	content, _ := brightdata.TruncateUTF8(string(body), brightdata.DefaultMaxChars())
	return scrapeResult(meta, content)
}

func handleDiscover(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	req := brightdata.DiscoverRequest{Query: strArg(m, "query"), Intent: strArg(m, "intent"),
		Country: strArg(m, "country"), City: strArg(m, "city"), Language: strArg(m, "language"),
		StartDate: strArg(m, "start_date"), EndDate: strArg(m, "end_date")}
	if n := intArg(m, "num_results", 0); n > 0 {
		req.NumResults = n
	}
	if kw, ok := m["filter_keywords"].([]interface{}); ok {
		for _, k := range kw {
			if s, ok := k.(string); ok && s != "" {
				req.FilterKeywords = append(req.FilterKeywords, s)
			}
		}
	}
	if v, ok := m["remove_duplicates"].(bool); ok {
		req.RemoveDuplicates = &v
	}
	cctx, cancel := pollCtx()
	defer cancel()
	data, err := client.Discover(cctx, req)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(string(data))
}

func pollCtx() (context.Context, context.CancelFunc) {
	sec := brightdata.PollingTimeoutSec()
	if sec < 30 {
		sec = 30
	}
	return context.WithTimeout(context.Background(), time.Duration(sec+5)*time.Second)
}
