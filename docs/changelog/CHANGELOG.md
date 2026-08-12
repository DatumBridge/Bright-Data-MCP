# Changelog

## 2026-08-12

### Added

- Full Bright Data MCP tool parity (60+ tools) per https://docs.brightdata.com/ai/mcp-server/tools
- Rapid tools: `search_engine`, `scrape_as_markdown`, `discover`, `scrape_as_html`, `scrape_batch`, `search_engine_batch`, `extract`, `session_stats`
- 50 `web_data_*` dataset tools via `datasets/v3/trigger` + snapshot polling
- `list_dataset_fields`, `search_dataset` for searchable LinkedIn datasets
- 13 `scraping_browser_*` tools (go-rod + Bright Data CDP)
- `BRIGHTDATA_PRO_MODE`, `BRIGHTDATA_GROUPS`, `BRIGHTDATA_TOOLS` tool filtering
- Embedded dataset catalog from official `@brightdata/mcp`

### Changed

- Default registration is Rapid mode (3 core + legacy aliases + health/stats)
- `scrape_url` tool args aligned with Web Unlocker API (`format=raw|json`, optional `data_format`)

### Fixed

- `scrape_url` appeared to return only `[UNTRUSTED_WEB_CONTENT]` when large HTML was JSON-escaped inside `content`
- Removed non-API `html`/`markdown` values from `format` (legacy aliases `html`→`raw`, `markdown`→`raw`+`data_format=markdown` still accepted)

### Removed

- N/A
