# Business Rules

1. Tool catalog follows [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) with DatumBridge Rapid/Pro/Groups filtering (`BRIGHTDATA_PRO_MODE`, `BRIGHTDATA_GROUPS`, `BRIGHTDATA_TOOLS`).
2. Default Rapid mode exposes search/scrape batch tools + legacy aliases (`scrape_url`, `search_serp`) and `bright_data_health` — not Pro-only tools (`session_stats`, `extract`, `web_data_*`, browser).
3. Unlocker/SERP/Browser zones come from credentials or env — never tool args.
4. Default scrape `format` is `raw` (Bright Data Web Unlocker API). Use `data_format=markdown` for LLM-friendly text.
5. `bright_data_health` must not call Bright Data unless `probe=true`.
6. Missing credentials/zones → MCP tool `isError` with actionable message (no secret echo).
7. SERP results are capped (`max_results`, default 10) on legacy `search_serp`.
8. **`discover` is not offered** — Discover API is retired (HTTP 410). Agents must use `search_engine`.
9. Operators are responsible for ToS/legal compliance of scrape targets.
10. Hosted Bright Data MCP (`mcp.brightdata.com`) is optional for external clients; Weaver Tool Registry uses this Go MCP (`mcpServer=bright-data`).
