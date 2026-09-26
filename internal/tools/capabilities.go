package tools

import "strings"

// capabilitiesForTool returns DatumBridge Tool Registry tags for tools/list `_meta`.
func capabilitiesForTool(name string) []string {
	switch {
	case strings.HasPrefix(name, "web_data_"):
		caps := []string{name, "data_fetching"}
		switch {
		case strings.Contains(name, "linkedin"), strings.Contains(name, "instagram"),
			strings.Contains(name, "facebook"), strings.Contains(name, "tiktok"),
			strings.Contains(name, "youtube"), strings.Contains(name, "reddit"),
			strings.Contains(name, "_x_"):
			caps = append(caps, "social_media")
		case strings.Contains(name, "amazon"), strings.Contains(name, "walmart"),
			strings.Contains(name, "ebay"), strings.Contains(name, "shop"),
			strings.Contains(name, "etsy"), strings.Contains(name, "bestbuy"):
			caps = append(caps, "ecommerce")
		}
		return caps
	case strings.HasPrefix(name, "scraping_browser_"):
		return []string{name, "browser_automation"}
	}

	switch name {
	case "search_engine", "search_serp", "search_engine_batch":
		return []string{name, "search"}
	case "scrape_as_markdown", "scrape_as_html", "scrape_batch", "scrape_url":
		return []string{name, "web_scraping"}
	case "extract":
		return []string{name, "data_processing"}
	case "discover":
		return []string{name, "search"}
	case "session_stats", "bright_data_health":
		return []string{name, "observability"}
	case "list_dataset_fields":
		return []string{name, "data_fetching"}
	case "search_dataset":
		return []string{name, "search"}
	default:
		return []string{name}
	}
}
