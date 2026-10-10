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

func baseProps(extra map[string]interface{}) map[string]interface{} {
	props := map[string]interface{}{
		"credentials_json": map[string]interface{}{
			"type":        "string",
			"description": "Vault JSON: {api_key, unlocker_zone, serp_zone, browser_zone}",
		},
		"credentials_path": map[string]interface{}{
			"type":        "string",
			"description": "Optional path under BRIGHTDATA_CREDENTIALS_DIR",
		},
	}
	for k, v := range extra {
		props[k] = v
	}
	return props
}

func schema(props map[string]interface{}, required []string) map[string]interface{} {
	s := map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

func clientFrom(raw json.RawMessage) (*brightdata.Client, map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		m = map[string]interface{}{}
	}
	credsJSON, _ := m["credentials_json"].(string)
	credsPath, _ := m["credentials_path"].(string)
	creds, err := brightdata.ParseCredentials(credsJSON, credsPath)
	if err != nil {
		return nil, m, err
	}
	return brightdata.NewClient(creds), m, nil
}

func boolArg(m map[string]interface{}, key string, def bool) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return def
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(t, "true") || t == "1"
	default:
		return def
	}
}

func intArg(m map[string]interface{}, key string, def int) int {
	v, ok := m[key]
	if !ok || v == nil {
		return def
	}
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	case string:
		var i int
		_, _ = fmt.Sscanf(t, "%d", &i)
		if i != 0 || t == "0" {
			return i
		}
	}
	return def
}

func strArg(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 95*time.Second)
}

func jsonResult(v interface{}) map[string]interface{} {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return mcp.ToolResultText(string(b))
}

func scrapeResult(meta map[string]interface{}, body string) map[string]interface{} {
	metaJSON, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	text := string(metaJSON) + "\n\n--- content ---\n\n" + body
	if ok, _ := meta["success"].(bool); !ok {
		return mcp.ToolResultError(text)
	}
	return mcp.ToolResultText(text)
}

// brightDataRefusal is the operator-facing reason when Unlocker returned no usable page.
// An empty body or x-brd-error is a failed search, not a successful empty document.
func brightDataRefusal(zone, target, format, dataFormat string, resp brightdata.DirectResponse) string {
	content := strings.TrimSpace(string(resp.Body))
	headers := brightdata.DiagnosticHeaders(resp.Header)
	msg := headers["x-brd-error"]
	if msg == "" {
		msg = headers["x-luminati-error"]
	}
	if msg == "" && content == "" {
		msg = "Bright Data returned an empty body"
	}
	if msg == "" {
		return ""
	}
	meta := unlockerResultMeta(zone, target, format, dataFormat, resp, content, false)
	encoded, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return msg
	}
	return msg + "\n" + string(encoded)
}

// unlockerResultMeta records the Bright Data request and the decoded response.
// The page body is returned separately and is not rewritten.
func unlockerResultMeta(zone, target, format, dataFormat string, resp brightdata.DirectResponse, content string, truncated bool) map[string]interface{} {
	headers := brightdata.DiagnosticHeaders(resp.Header)
	meta := map[string]interface{}{
		"success":            resp.Status >= 200 && resp.Status < 300,
		"url":                target,
		"status":             resp.Status,
		"zone":               zone,
		"format":             format,
		"truncated":          truncated,
		"char_count":         len([]rune(content)),
		"empty_body":         strings.TrimSpace(content) == "",
		"brightdata_headers": headers,
	}
	if dataFormat != "" {
		meta["data_format"] = dataFormat
	}
	if msg := headers["x-brd-error"]; msg != "" {
		meta["brightdata_error"] = msg
		meta["success"] = false
	} else if msg := headers["x-luminati-error"]; msg != "" {
		meta["brightdata_error"] = msg
		meta["success"] = false
	}
	if strings.TrimSpace(content) == "" {
		meta["success"] = false
		if _, ok := meta["brightdata_error"]; !ok {
			meta["brightdata_error"] = "Bright Data returned an empty body"
		}
	}
	return meta
}
