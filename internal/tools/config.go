package tools

import (
	"os"
	"strings"
)

// ServerConfig mirrors Bright Data MCP PRO_MODE / GROUPS / TOOLS env filtering.
type ServerConfig struct {
	ProMode     bool
	Groups      []string
	CustomTools []string
}

func LoadServerConfig() ServerConfig {
	cfg := ServerConfig{
		ProMode: strings.EqualFold(os.Getenv("BRIGHTDATA_PRO_MODE"), "true") ||
			strings.EqualFold(os.Getenv("PRO_MODE"), "true"),
	}
	if v := strings.TrimSpace(os.Getenv("BRIGHTDATA_GROUPS")); v != "" {
		cfg.Groups = splitCSV(v)
	} else if v := strings.TrimSpace(os.Getenv("GROUPS")); v != "" {
		cfg.Groups = splitCSV(v)
	}
	if v := strings.TrimSpace(os.Getenv("BRIGHTDATA_TOOLS")); v != "" {
		cfg.CustomTools = splitCSV(v)
	} else if v := strings.TrimSpace(os.Getenv("TOOLS")); v != "" {
		cfg.CustomTools = splitCSV(v)
	}
	return cfg
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, strings.ToLower(p))
		}
	}
	return out
}

// rapidTools are enabled by default (Rapid/Free mode).
var rapidTools = map[string]bool{
	"search_engine":       true,
	"scrape_as_markdown":  true,
	"discover":            true,
	"bright_data_health":  true,
	"session_stats":       true,
	// Legacy DatumBridge aliases
	"scrape_url":  true,
	"search_serp": true,
}

func (c ServerConfig) toolEnabled(name string) bool {
	if c.ProMode {
		return true
	}
	if len(c.CustomTools) > 0 {
		for _, t := range c.CustomTools {
			if t == strings.ToLower(name) {
				return true
			}
		}
		return false
	}
	if len(c.Groups) > 0 {
		allowed := buildAllowedFromGroups(c.Groups)
		return allowed[strings.ToLower(name)]
	}
	return rapidTools[name]
}

func buildAllowedFromGroups(groups []string) map[string]bool {
	allowed := map[string]bool{}
	for _, g := range groups {
		for _, t := range toolsInGroup(g) {
			allowed[t] = true
		}
	}
	return allowed
}
