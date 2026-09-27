package brightdata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// BuildEngineSearchURL builds Google/Bing/Yandex search URLs (official MCP parity).
func BuildEngineSearchURL(engine, query, cursor, geo string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	engine = strings.ToLower(strings.TrimSpace(engine))
	if engine == "" {
		engine = "google"
	}
	page := 0
	if cursor != "" {
		if p, err := strconv.Atoi(cursor); err == nil && p >= 0 {
			page = p
		}
	}
	q := url.QueryEscape(query)
	start := page * 10
	switch engine {
	case "yandex":
		return fmt.Sprintf("https://yandex.com/search/?text=%s&p=%d", q, page), nil
	case "bing":
		return fmt.Sprintf("https://www.bing.com/search?q=%s&first=%d", q, start+1), nil
	case "google":
		gl := ""
		if geo != "" {
			gl = "&gl=" + url.QueryEscape(strings.ToLower(geo))
		}
		return fmt.Sprintf("https://www.google.com/search?q=%s&start=%d%s", q, start, gl), nil
	default:
		return "", fmt.Errorf("unsupported engine %q (use google, bing, or yandex)", engine)
	}
}

// SearchQuery is one item of search_engine_batch.
type SearchQuery struct {
	Query  string
	Engine string
	Cursor string
	Geo    string
}

// NormalizeBatchQueries accepts a list of strings, objects with a query field,
// or a JSON string of that list. The planner often sends queries as strings.
func NormalizeBatchQueries(raw interface{}) ([]SearchQuery, error) {
	if raw == nil {
		return nil, fmt.Errorf("queries array is required (max 10)")
	}
	if text, ok := raw.(string); ok {
		text = strings.TrimSpace(text)
		if text == "" {
			return nil, fmt.Errorf("queries array is required (max 10)")
		}
		var decoded interface{}
		if err := json.Unmarshal([]byte(text), &decoded); err != nil {
			return []SearchQuery{{Query: text}}, nil
		}
		return NormalizeBatchQueries(decoded)
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil, fmt.Errorf("queries array is required (max 10)")
	}
	if len(list) > 10 {
		list = list[:10]
	}
	out := make([]SearchQuery, 0, len(list))
	for _, item := range list {
		switch v := item.(type) {
		case string:
			q := strings.TrimSpace(v)
			if q == "" {
				return nil, fmt.Errorf("query is required")
			}
			out = append(out, SearchQuery{Query: q})
		case map[string]interface{}:
			q := mapString(v, "query", "q", "text")
			if q == "" || q == "<nil>" {
				return nil, fmt.Errorf("query is required")
			}
			out = append(out, SearchQuery{
				Query:  q,
				Engine: mapString(v, "engine"),
				Cursor: mapString(v, "cursor"),
				Geo:    mapString(v, "geo_location", "geo"),
			})
		default:
			return nil, fmt.Errorf("query is required")
		}
	}
	return out, nil
}

func mapString(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		s, isString := v.(string)
		if !isString {
			s = fmt.Sprint(v)
		}
		s = strings.TrimSpace(s)
		if s == "" || s == "<nil>" {
			continue
		}
		return s
	}
	return ""
}

// ParseGoogleSearchResponse parses Bright Data parsed_light / brd_json Google SERP.
func ParseGoogleSearchResponse(body []byte) (map[string]interface{}, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		snippet := string(body)
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		return nil, fmt.Errorf("unexpected non-JSON Google SERP response: %s", snippet)
	}
	organicRaw, _ := raw["organic"].([]interface{})
	clean := make([]map[string]string, 0, len(organicRaw))
	for _, item := range organicRaw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		link, _ := m["link"].(string)
		title, _ := m["title"].(string)
		desc, _ := m["description"].(string)
		link = strings.TrimSpace(link)
		title = strings.TrimSpace(title)
		desc = strings.TrimSpace(desc)
		if link == "" || title == "" {
			continue
		}
		clean = append(clean, map[string]string{
			"link": link, "title": title, "description": desc,
		})
	}
	return map[string]interface{}{"organic": clean}, nil
}
