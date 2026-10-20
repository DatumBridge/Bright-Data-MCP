# Workflows

## Local connect

1. Copy a Bright Data API token. Lab uses Direct API zones `web_unlocker2` (Unlocker) and `serp_api2` (legacy SERP). They are the process defaults when env/vault omit the names.
2. Set `BRIGHTDATA_API_KEY` in `.env` or Tool Registry secret vars. Zone public vars are optional on lab.
3. `go run ./cmd/api` or Tool Registry Setup/Publish (port `8011`). Do not set `HTTP_PROXY` to Superproxy on the MCP pod.
4. `initialize` → `tools/list` → `tools/call` / `bright_data_health`

## Rapid web search

1. Agent calls `search_engine` with `query` (optional `engine`, `cursor`, `geo_location`)
2. Server builds the Google/Bing/Yandex URL and calls Unlocker `/request` with `data_format=markdown`
3. The page text is returned as soon as Bright Data fetches it. JSON organic hits are kept when the body is already JSON
4. For many queries, use `search_engine_batch` (up to 10). `search_serp` uses the same fast fetch unless `brd_json=true`

## Scrape page (Markdown)

1. Agent calls `scrape_as_markdown` with `url` (or batch via `scrape_batch`)
2. Server uses Unlocker zone + Bearer key; `format=raw`, `data_format=markdown`
3. Returns metadata (zone, HTTP status, `brightdata_headers`) and the decoded Bright Data body after `--- content ---`. An empty body or `x-brd-error` is `isError`. The text does not include an untrusted-content banner.

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
