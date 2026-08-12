package tools

// toolsInGroup returns tool names for a Bright Data MCP group id (official tool_groups.js parity).
func toolsInGroup(groupID string) []string {
	base := []string{"search_engine", "scrape_as_markdown", "discover"}
	switch groupID {
	case "ecommerce":
		return append(base,
			"web_data_amazon_product", "web_data_amazon_product_reviews", "web_data_amazon_product_search",
			"web_data_walmart_product", "web_data_walmart_seller", "web_data_ebay_product",
			"web_data_homedepot_products", "web_data_zara_products", "web_data_etsy_products",
			"web_data_bestbuy_products", "web_data_google_shopping",
		)
	case "social":
		return append(base,
			"web_data_linkedin_person_profile", "web_data_linkedin_company_profile",
			"web_data_linkedin_job_listings", "web_data_linkedin_posts", "web_data_linkedin_people_search",
			"list_dataset_fields", "search_dataset",
			"web_data_instagram_profiles", "web_data_instagram_posts", "web_data_instagram_reels",
			"web_data_instagram_comments", "web_data_facebook_posts", "web_data_facebook_marketplace_listings",
			"web_data_facebook_company_reviews", "web_data_facebook_events",
			"web_data_tiktok_profiles", "web_data_tiktok_posts", "web_data_tiktok_shop", "web_data_tiktok_comments",
			"web_data_x_posts", "web_data_x_profile_posts",
			"web_data_youtube_profiles", "web_data_youtube_comments", "web_data_youtube_videos",
			"web_data_reddit_posts", "web_data_reddit_comments",
		)
	case "browser":
		return append(base,
			"scraping_browser_navigate", "scraping_browser_go_back", "scraping_browser_go_forward",
			"scraping_browser_snapshot", "scraping_browser_click_ref", "scraping_browser_type_ref",
			"scraping_browser_screenshot", "scraping_browser_network_requests",
			"scraping_browser_wait_for_ref", "scraping_browser_get_text", "scraping_browser_get_html",
			"scraping_browser_scroll", "scraping_browser_scroll_to_ref",
		)
	case "finance":
		return append(base, "web_data_yahoo_finance_business")
	case "business":
		return append(base,
			"web_data_crunchbase_company", "web_data_zoominfo_company_profile",
			"web_data_google_maps_reviews", "web_data_zillow_properties_listing",
			"web_data_booking_hotel_listings", "list_dataset_fields", "search_dataset",
		)
	case "research":
		return append(base, "web_data_github_repository_file", "web_data_reuter_news")
	case "app_stores":
		return append(base, "web_data_google_play_store", "web_data_apple_app_store")
	case "travel":
		return append(base, "web_data_booking_hotel_listings")
	case "geo":
		return append(base,
			"web_data_chatgpt_ai_insights", "web_data_grok_ai_insights", "web_data_perplexity_ai_insights",
		)
	case "code":
		return append(base, "web_data_npm_package", "web_data_pypi_package")
	case "advanced_scraping":
		return append(base,
			"search_engine_batch", "scrape_batch", "scrape_as_html", "extract", "session_stats",
		)
	case "all":
		return allToolNames()
	default:
		return nil
	}
}

func allToolNames() []string {
	names := []string{
		"search_engine", "scrape_as_markdown", "discover",
		"search_engine_batch", "scrape_batch", "scrape_as_html", "extract", "session_stats",
		"list_dataset_fields", "search_dataset", "bright_data_health",
		"scraping_browser_navigate", "scraping_browser_go_back", "scraping_browser_go_forward",
		"scraping_browser_snapshot", "scraping_browser_click_ref", "scraping_browser_type_ref",
		"scraping_browser_screenshot", "scraping_browser_network_requests",
		"scraping_browser_wait_for_ref", "scraping_browser_get_text", "scraping_browser_get_html",
		"scraping_browser_scroll", "scraping_browser_scroll_to_ref",
		"scrape_url", "search_serp",
	}
	for _, ds := range datasetToolNames() {
		names = append(names, ds)
	}
	return names
}

func datasetToolNames() []string {
	catalog, err := loadDatasetCatalogCached()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(catalog))
	for _, t := range catalog {
		out = append(out, "web_data_"+t.ID)
	}
	return out
}
