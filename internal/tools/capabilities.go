package tools

import "strings"

// capabilitiesForTool returns DatumBridge Tool Registry tags for tools/list `_meta`.
func capabilitiesForTool(name string) []string {
	switch {
	case strings.HasPrefix(name, "web_data_"):
		caps := []string{"web_scraping", "data_fetching", "bright_data"}
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
		return []string{"browser_automation", "web_scraping", "bright_data"}
	}

	switch name {
	case "search_engine", "search_serp", "search_engine_batch":
		return []string{"search", "web_scraping", "bright_data"}
	case "scrape_as_markdown", "scrape_as_html", "scrape_batch", "scrape_url":
		return []string{"web_scraping", "data_fetching", "bright_data"}
	case "extract":
		return []string{"web_scraping", "data_processing", "bright_data"}
	case "discover":
		return []string{"search", "web_scraping", "discover", "bright_data"}
	case "session_stats", "bright_data_health":
		return []string{"web_scraping", "observability", "bright_data"}
	case "list_dataset_fields":
		return []string{"web_scraping", "data_fetching", "bright_data"}
	case "search_dataset":
		return []string{"web_scraping", "search", "bright_data"}
	default:
		return []string{"web_scraping", "bright_data"}
	}
}
