# Changelog

## 2026-09-05

### Changed

- Rapid default tools aligned with official Bright Data MCP: `search_engine`, `scrape_as_markdown`, `search_engine_batch`, `scrape_batch` (+ DatumBridge aliases `scrape_url` / `search_serp` / `bright_data_health`)
- `session_stats` moved to Pro / `advanced_scraping` only (no longer Rapid)
- Group base tools are `search_engine` + `scrape_as_markdown` only
- `business` group matches official catalog (booking stays in `travel` only)
- `research` group is GitHub file only; `web_data_reuter_news` / `web_data_reddit_comments` remain Pro extras via `datasets.json`
- Catalog reference updated to https://docs.brightdata.com/products/mcp-server/tools

### Removed

- **`discover` tool** — Bright Data Discover API returns `HTTP 410: Discover API is no longer available`. Use `search_engine` / `search_engine_batch` instead. Deleted `internal/brightdata/discover.go` and MCP registration.

### Added

- Docs note: hosted `https://mcp.brightdata.com/sse?token=…` is optional; Weaver uses this Go HTTP MCP

## 2026-08-12

### Added

- Rapid tools at introduction: `search_engine`, `scrape_as_markdown`, plus Pro tools `scrape_as_html`, `scrape_batch`, `search_engine_batch`, `extract`, `session_stats` (Rapid/Pro split refined on 2026-09-05)
- 50 `web_data_*` dataset tools via `datasets/v3/trigger` + snapshot polling
- `list_dataset_fields`, `search_dataset` for searchable LinkedIn datasets
- 13 `scraping_browser_*` tools (go-rod + Bright Data CDP)
- `BRIGHTDATA_PRO_MODE`, `BRIGHTDATA_GROUPS`, `BRIGHTDATA_TOOLS` tool filtering
- Embedded dataset catalog from official `@brightdata/mcp`

### Changed

- Default registration is Rapid mode (core search/scrape + legacy aliases + health)
- `scrape_url` tool args aligned with Web Unlocker API (`format=raw|json`, optional `data_format`)

### Fixed

- Docker/deploy build failed: `go:embed data/datasets.json` missing because `.gitignore` rule `data/` excluded `internal/brightdata/data/` from git clone

### Removed

- N/A (see 2026-09-05 for `discover` removal)
