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
		"Scrape search results from Google, Bing, or Yandex as Markdown. JSON organic results are returned only when Bright Data already sent JSON.",
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
		"Run up to 10 search queries in one request and return Markdown. queries is an array of strings, or objects with a query field. JSON organic hits are kept when Bright Data already sent JSON.",
		baseProps(map[string]interface{}{
			"queries": map[string]interface{}{
				"type":        "array",
				"description": "Search strings, or objects {query, engine, cursor, geo_location}.",
				"items": map[string]interface{}{
					"oneOf": []interface{}{
						map[string]interface{}{"type": "string"},
						map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"query":        map[string]interface{}{"type": "string"},
								"engine":       map[string]interface{}{"type": "string"},
								"cursor":       map[string]interface{}{"type": "string"},
								"geo_location": map[string]interface{}{"type": "string"},
							},
							"required": []string{"query"},
						},
					},
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
		"[Legacy alias] Use search_engine. SERP via Web Unlocker. brd_json defaults to false so the page is returned without Bright Data JSON parsing. Set brd_json=true to request parsed JSON.",
		baseProps(map[string]interface{}{
			"query":       map[string]interface{}{"type": "string"},
			"engine":      map[string]interface{}{"type": "string", "enum": []string{"google", "bing"}},
			"country":     map[string]interface{}{"type": "string"},
			"language":    map[string]interface{}{"type": "string"},
			"start":       map[string]interface{}{"type": "integer"},
			"brd_json":    map[string]interface{}{"type": "boolean", "default": false, "description": "Request Bright Data parsed JSON. Slower. Default false returns the page as text."},
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
	opts := brightdata.RequestOpts{Zone: zone, URL: searchURL, Format: "raw", DataFormat: "markdown"}
	cctx, cancel := pollCtx()
	defer cancel()
	resp, err := client.Request(cctx, opts)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return searchBodyResult(resp.Body)
}

// searchBodyResult keeps structured organic hits when Bright Data sent JSON.
// Markdown or HTML is returned as text so a non-JSON page is still usable.
func searchBodyResult(body []byte) map[string]interface{} {
	if parsed, err := brightdata.ParseGoogleSearchResponse(body); err == nil {
		if organic, ok := parsed["organic"].([]map[string]string); ok && len(organic) > 0 {
			return jsonResult(parsed)
		}
	}
	text, _ := brightdata.TruncateUTF8(string(body), brightdata.DefaultMaxChars())
	return untrustedTextResult(text)
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
	resp, err := client.Request(cctx, opts)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	content, truncated := brightdata.TruncateUTF8(string(resp.Body), brightdata.DefaultMaxChars())
	return scrapeResult(unlockerResultMeta(zone, target, "raw", dataFormat, resp, content, truncated), content)
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
	queries, err := brightdata.NormalizeBatchQueries(m["queries"])
	if err != nil {
		if q := strArg(m, "query"); q != "" {
			queries = []brightdata.SearchQuery{{
				Query:  q,
				Engine: strArg(m, "engine"),
				Cursor: strArg(m, "cursor"),
				Geo:    strArg(m, "geo_location"),
			}}
		} else {
			return mcp.ToolResultError(err.Error())
		}
	}
	type result struct {
		Query  string      `json:"query"`
		Engine string      `json:"engine"`
		Result interface{} `json:"result,omitempty"`
		Error  string      `json:"error,omitempty"`
	}
	out := make([]result, 0, len(queries))
	failures := 0
	for _, item := range queries {
		query := item.Query
		engine := item.Engine
		if engine == "" || engine == "<nil>" {
			engine = "google"
		}
		searchURL, err := brightdata.BuildEngineSearchURL(engine, query, item.Cursor, item.Geo)
		if err != nil {
			failures++
			out = append(out, result{Query: query, Engine: engine, Error: err.Error()})
			continue
		}
		opts := brightdata.RequestOpts{Zone: zone, URL: searchURL, Format: "raw", DataFormat: "markdown"}
		cctx, cancel := pollCtx()
		resp, reqErr := client.Request(cctx, opts)
		cancel()
		body := resp.Body
		if reqErr != nil {
			failures++
			out = append(out, result{Query: query, Engine: engine, Error: reqErr.Error()})
			continue
		}
		if parsed, perr := brightdata.ParseGoogleSearchResponse(body); perr == nil {
			if organic, ok := parsed["organic"].([]map[string]string); ok && len(organic) > 0 {
				out = append(out, result{Query: query, Engine: engine, Result: parsed})
				continue
			}
		}
		text, _ := brightdata.TruncateUTF8(string(body), brightdata.DefaultMaxChars())
		out = append(out, result{Query: query, Engine: engine, Result: text})
	}
	if failures == len(out) {
		payload, _ := json.Marshal(out)
		return mcp.ToolResultError(string(payload))
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
		URL       string            `json:"url"`
		Status    int               `json:"status,omitempty"`
		Zone      string            `json:"zone,omitempty"`
		Headers   map[string]string `json:"brightdata_headers,omitempty"`
		CharCount int               `json:"char_count,omitempty"`
		EmptyBody bool              `json:"empty_body,omitempty"`
		Content   string            `json:"content,omitempty"`
		Error     string            `json:"error,omitempty"`
	}
	out := make([]pair, 0, len(urlsRaw))
	for _, u := range urlsRaw {
		target, _ := u.(string)
		target = strings.TrimSpace(target)
		opts := brightdata.RequestOpts{Zone: zone, URL: target, Format: "raw", DataFormat: "markdown"}
		cctx, cancel := pollCtx()
		resp, reqErr := client.Request(cctx, opts)
		cancel()
		headers := brightdata.DiagnosticHeaders(resp.Header)
		if reqErr != nil {
			item := pair{URL: target, Status: resp.Status, Zone: zone, Headers: headers, Error: reqErr.Error()}
			out = append(out, item)
			continue
		}
		content, _ := brightdata.TruncateUTF8(string(resp.Body), brightdata.DefaultMaxChars())
		out = append(out, pair{
			URL:       target,
			Status:    resp.Status,
			Zone:      zone,
			Headers:   headers,
			CharCount: len([]rune(content)),
			EmptyBody: strings.TrimSpace(content) == "",
			Content:   content,
		})
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
	resp, err := client.Request(cctx, brightdata.RequestOpts{
		Zone: zone, URL: target, Format: "raw", DataFormat: "markdown",
	})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	content, truncated := brightdata.TruncateUTF8(string(resp.Body), brightdata.DefaultMaxChars())
	meta := unlockerResultMeta(zone, target, "raw", "markdown", resp, content, truncated)
	meta["note"] = "Structure the markdown below as JSON using your LLM per extraction_prompt"
	if prompt != "" {
		meta["extraction_prompt"] = prompt
	}
	return scrapeResult(meta, content)
}

func pollCtx() (context.Context, context.CancelFunc) {
	sec := brightdata.PollingTimeoutSec()
	if sec < 30 {
		sec = 30
	}
	return context.WithTimeout(context.Background(), time.Duration(sec+5)*time.Second)
}
