package tools

import "testing"

func TestCapabilitiesForTool_distinctFamilies(t *testing.T) {
	search := capabilitiesForTool("search_engine")
	scrape := capabilitiesForTool("scrape_as_markdown")
	browser := capabilitiesForTool("scraping_browser_navigate")
	social := capabilitiesForTool("web_data_linkedin_posts")
	if equalStrings(search, scrape) {
		t.Fatalf("search and scrape should differ: %v vs %v", search, scrape)
	}
	if !contains(browser, "browser_automation") {
		t.Fatalf("browser tool missing browser_automation: %v", browser)
	}
	if !contains(social, "social_media") {
		t.Fatalf("linkedin web_data missing social_media: %v", social)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func contains(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}
