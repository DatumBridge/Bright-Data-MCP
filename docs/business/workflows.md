# Workflows

## Local connect

1. Create Unlocker (+ Browser for Pro) zones; copy API key
2. Set `.env` from `.env.example` (`BRIGHTDATA_PRO_MODE=true` for full catalog)
3. `go run ./cmd/api`
4. `initialize` → `tools/list` → `tools/call`

## Rapid web search

1. Agent calls `search_engine` with `query` (optional `engine`, `cursor`, `geo_location`)
2. Server builds Google/Bing/Yandex URL and calls Unlocker `/request`
3. Google returns structured organic JSON; Bing/Yandex return Markdown
4. For many queries, use `search_engine_batch` (up to 10)

## Scrape page (Markdown)

1. Agent calls `scrape_as_markdown` with `url` (or batch via `scrape_batch`)
2. Server uses Unlocker zone + Bearer key; `format=raw`, `data_format=markdown`
3. Returns untrusted-prefixed content

## Legacy aliases

1. `scrape_url` / `search_serp` remain for older agents; prefer official names above
2. `bright_data_health` validates config; `probe=true` is billable

## Pro structured / browser

1. Enable `BRIGHTDATA_PRO_MODE=true` (or `BRIGHTDATA_GROUPS=…`)
2. `web_data_*` → Datasets trigger + snapshot poll
3. `scraping_browser_*` → Scraping Browser CDP (requires `BRIGHTDATA_BROWSER_ZONE`)

## Migration: discover removed

1. Do not call `discover` (unknown tool / previously HTTP 410)
2. Use `search_engine` or `search_engine_batch` instead
3. After deploy, republish Weaver Tool Registry for `bright-data` so `tools/list` drops `discover`
