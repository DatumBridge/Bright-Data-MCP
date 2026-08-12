# Business Rules

1. Only three tools in v1: `scrape_url`, `search_serp`, `bright_data_health`.
2. Unlocker/SERP zones come from credentials or env — never tool args.
3. Default scrape `format` is `raw` (Bright Data Web Unlocker API). Use `data_format=markdown` for LLM-friendly text.
4. `bright_data_health` must not call Bright Data unless `probe=true`.
5. Missing credentials/zones → MCP tool `isError` with actionable message (no secret echo).
6. SERP results are capped (`max_results`, default 10).
7. Operators are responsible for ToS/legal compliance of scrape targets.
