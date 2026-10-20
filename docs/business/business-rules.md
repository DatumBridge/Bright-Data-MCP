# Business Rules

1. Tool catalog follows [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) with DatumBridge Rapid/Pro/Groups filtering (`BRIGHTDATA_PRO_MODE`, `BRIGHTDATA_GROUPS`, `BRIGHTDATA_TOOLS`).
2. Default Rapid mode exposes search/scrape batch tools + legacy aliases (`scrape_url`, `search_serp`) and `bright_data_health` — not Pro-only tools (`session_stats`, `extract`, `web_data_*`, browser).
3. Unlocker/SERP/Browser zones come from credentials or env — never tool args. Lab defaults when omitted: Unlocker `web_unlocker2`, SERP `serp_api2`. Browser zone has no default.
4. Default scrape `format` is `raw` (Bright Data Web Unlocker API). Use `data_format=markdown` for LLM-friendly text.
5. `bright_data_health` must not call Bright Data unless `probe=true`.
6. Missing API key → MCP tool `isError` with an actionable message (no secret echo). Missing zone names use the lab defaults above.
7. SERP results are capped (`max_results`, default 10) on legacy `search_serp`.
8. **`discover` is not offered** — Discover API is retired (HTTP 410). Agents must use `search_engine`.
9. Operators are responsible for ToS/legal compliance of scrape targets.
10. Hosted Bright Data MCP (`mcp.brightdata.com`) is optional for external clients; Weaver Tool Registry uses this Go MCP (`mcpServer=bright-data`).
11. Unlocker tool output must show the zone that was called and Bright Data diagnostic headers, and must place the decoded response body after `--- content ---` without rewriting it. An empty HTTP 200 or `x-brd-error` / `x-luminati-error` is `isError: true`, with `empty_body` and the headers in that error text, so Lab and Production can be compared. The tool text does not include an `[UNTRUSTED_WEB_CONTENT]` banner.
12. Search tools return the fetched page as text by default. An empty page or a Bright Data error header is `isError`, not a successful document whose only text is a warning. Parsed Google JSON (`brd_json` / `parsed_light`) is opt-in on `search_serp` because that parse holds the request until Bright Data finishes structuring the SERP.
