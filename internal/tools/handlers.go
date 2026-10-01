package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
	"github.com/datumbridge/bright-data-mcp/internal/mcp"
)

func handleScrapeURL(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone := client.UnlockerZone()
	if zone == "" {
		return mcp.ToolResultError("unlocker_zone required (credentials_json.unlocker_zone or BRIGHTDATA_UNLOCKER_ZONE)")
	}
	target := strArg(m, "url")
	format, dataFormat, err := parseUnlockerFormats(m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	country := strArg(m, "country")
	if country != "" && len(country) != 2 {
		return mcp.ToolResultError("country must be a 2-letter ISO code")
	}
	maxChars := intArg(m, "max_chars", brightdata.DefaultMaxChars())
	if maxChars < 1 {
		maxChars = brightdata.DefaultMaxChars()
	}
	if maxChars > 500_000 {
		maxChars = 500_000
	}

	opts := brightdata.RequestOpts{
		Zone:       zone,
		URL:        target,
		Format:     format,
		DataFormat: dataFormat,
		Country:    country,
	}

	cctx, cancel := ctx()
	defer cancel()
	resp, err := client.Request(cctx, opts)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	content, truncated := brightdata.TruncateUTF8(string(resp.Body), maxChars)
	meta := unlockerResultMeta(zone, target, format, dataFormat, resp, content, truncated)
	if country != "" {
		meta["country"] = strings.ToLower(country)
	}
	return scrapeResult(meta, content)
}

// parseUnlockerFormats maps tool args to Bright Data Web Unlocker API fields.
// API format is only "raw" or "json". Optional data_format (markdown|screenshot) applies when format=raw.
func parseUnlockerFormats(m map[string]interface{}) (format, dataFormat string, err error) {
	format = strings.ToLower(strArg(m, "format"))
	dataFormat = strings.ToLower(strArg(m, "data_format"))

	// Backward compatibility for pre-API-alignment tool args.
	switch format {
	case "", "html":
		format = "raw"
	case "markdown":
		format = "raw"
		if dataFormat == "" {
			dataFormat = "markdown"
		}
	}

	if format != "raw" && format != "json" {
		return "", "", fmt.Errorf("format must be raw or json (Bright Data Web Unlocker API)")
	}
	if dataFormat != "" {
		if format != "raw" {
			return "", "", fmt.Errorf("data_format is only valid when format=raw")
		}
		if dataFormat != "markdown" && dataFormat != "screenshot" {
			return "", "", fmt.Errorf("data_format must be markdown or screenshot")
		}
	}
	return format, dataFormat, nil
}

func handleSearchSERP(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	zone := client.SerpZone()
	if zone == "" {
		return mcp.ToolResultError("serp_zone required (credentials_json.serp_zone or BRIGHTDATA_SERP_ZONE)")
	}
	query := strArg(m, "query")
	engine := strArg(m, "engine")
	country := strArg(m, "country")
	language := strArg(m, "language")
	start := intArg(m, "start", 0)
	brdJSON := boolArg(m, "brd_json", true)
	maxResults := intArg(m, "max_results", 10)
	if maxResults < 1 {
		maxResults = 10
	}
	if maxResults > 50 {
		maxResults = 50
	}

	serpURL, err := brightdata.BuildSERPURL(query, engine, country, language, start, brdJSON)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if engine == "" {
		engine = "google"
	}

	cctx, cancel := ctx()
	defer cancel()
	resp, err := client.Request(cctx, brightdata.RequestOpts{
		Zone:   zone,
		URL:    serpURL,
		Format: "raw",
	})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	body := resp.Body
	status := resp.Status

	results, parseNote := parseSERPResults(body, maxResults)
	out := map[string]interface{}{
		"success":      true,
		"query":        query,
		"engine":       strings.ToLower(engine),
		"status":       status,
		"result_count": len(results),
		"results":      results,
	}
	if parseNote != "" {
		out["parse_note"] = parseNote
	}
	// If not JSON, return truncated raw text instead of dumping HTML wholesale into results.
	if len(results) == 0 && !json.Valid(body) {
		rawText, truncated := brightdata.TruncateUTF8(string(body), 20_000)
		out["raw_excerpt"] = rawText
		out["truncated"] = truncated
	}
	return jsonResult(out)
}

func handleHealth(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	probe := boolArg(m, "probe", false)
	out := map[string]interface{}{
		"success":         true,
		"api_key_masked":  client.APIKeyMasked(),
		"unlocker_zone":   nonemptyOrMissing(client.UnlockerZone()),
		"serp_zone":       nonemptyOrMissing(client.SerpZone()),
		"probe_requested": probe,
	}
	if !probe {
		out["message"] = "Config present. Set probe=true to run a live Unlocker request (billable)."
		b, _ := json.MarshalIndent(out, "", "  ")
		return mcp.ToolResultText(string(b))
	}
	if client.UnlockerZone() == "" {
		return mcp.ToolResultError("probe requires unlocker_zone")
	}
	cctx, cancel := ctx()
	defer cancel()
	resp, err := client.Request(cctx, brightdata.RequestOpts{
		Zone:   client.UnlockerZone(),
		URL:    "https://example.com",
		Format: "raw",
	})
	if err != nil {
		out["success"] = false
		out["probe_error"] = err.Error()
		out["probe_status"] = resp.Status
		out["brightdata_headers"] = brightdata.DiagnosticHeaders(resp.Header)
		b, _ := json.MarshalIndent(out, "", "  ")
		return mcp.ToolResultError(string(b))
	}
	out["probe_status"] = resp.Status
	out["brightdata_headers"] = brightdata.DiagnosticHeaders(resp.Header)
	out["message"] = "Live Unlocker probe succeeded"
	b, _ := json.MarshalIndent(out, "", "  ")
	return mcp.ToolResultText(string(b))
}

func nonemptyOrMissing(s string) string {
	if strings.TrimSpace(s) == "" {
		return "(missing)"
	}
	return s
}

func parseSERPResults(body []byte, maxResults int) ([]map[string]interface{}, string) {
	var root interface{}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, "response was not JSON; set brd_json=true or inspect raw_excerpt"
	}

	// Common Bright Data SERP JSON shapes: organic / organic_results / body.organic
	candidates := [][]interface{}{}
	switch t := root.(type) {
	case map[string]interface{}:
		for _, key := range []string{"organic", "organic_results", "results"} {
			if arr, ok := t[key].([]interface{}); ok {
				candidates = append(candidates, arr)
			}
		}
		if bodyMap, ok := t["body"].(map[string]interface{}); ok {
			for _, key := range []string{"organic", "organic_results", "results"} {
				if arr, ok := bodyMap[key].([]interface{}); ok {
					candidates = append(candidates, arr)
				}
			}
		}
	case []interface{}:
		candidates = append(candidates, t)
	}

	var src []interface{}
	if len(candidates) > 0 {
		src = candidates[0]
	}
	out := make([]map[string]interface{}, 0, maxResults)
	for _, item := range src {
		if len(out) >= maxResults {
			break
		}
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := map[string]interface{}{}
		if title := firstString(m, "title", "name"); title != "" {
			entry["title"] = title
		}
		if link := firstString(m, "link", "url", "href"); link != "" {
			entry["link"] = link
		}
		if snip := firstString(m, "snippet", "description", "desc"); snip != "" {
			entry["snippet"] = snip
		}
		if len(entry) == 0 {
			continue
		}
		out = append(out, entry)
	}
	if len(out) == 0 {
		return out, "could not extract organic results from SERP JSON"
	}
	return out, ""
}

func firstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}
