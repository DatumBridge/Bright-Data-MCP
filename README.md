# Bright Data MCP Server

Go MCP **tool-server** aligned with [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) (60+ tools).

**mcpServer id:** `bright-data` (aliases: `bright-data-mcp`)

## Tool modes

| Mode | Env | Tools |
|------|-----|-------|
| **Rapid (default)** | `BRIGHTDATA_PRO_MODE=false` | `search_engine`, `scrape_as_markdown`, `search_engine_batch`, `scrape_batch`, legacy `scrape_url` / `search_serp`, `bright_data_health` |
| **Pro** | `BRIGHTDATA_PRO_MODE=true` | All 60+ tools including `web_data_*`, browser automation, `extract`, `session_stats`, dataset search |
| **Groups** | `BRIGHTDATA_GROUPS=ecommerce,social` | Subset per [official groups](https://docs.brightdata.com/products/mcp-server/tools) |
| **Custom** | `BRIGHTDATA_TOOLS=tool1,tool2` | Explicit allowlist |

## Tool categories (Pro)

- **Search/scrape:** `search_engine`, `scrape_as_markdown`, `scrape_as_html`, `scrape_batch`, `search_engine_batch`, `extract`, `session_stats`
- **Structured data:** `web_data_amazon_product`, `web_data_linkedin_person_profile`, … (catalog in `internal/brightdata/data/datasets.json`)
- **Dataset search:** `list_dataset_fields`, `search_dataset`
- **Browser:** `scraping_browser_*` (13 tools via Scraping Browser + CDP)
- **GEO / code:** `web_data_chatgpt_ai_insights`, `web_data_npm_package`, …
- **Legacy aliases:** `scrape_url`, `search_serp`

## Deprecated / removed

- **`discover`** — Bright Data Discover API returns `HTTP 410` (`Discover API is no longer available`). Use `search_engine` (or `search_engine_batch`) instead.

## Setup

```bash
BRIGHTDATA_API_KEY=...
BRIGHTDATA_UNLOCKER_ZONE=web_unlocker1   # Web Unlocker / SERP via /request
BRIGHTDATA_BROWSER_ZONE=mcp_browser      # Pro browser tools
BRIGHTDATA_PRO_MODE=true                 # enable all tools
```

Vault `credentials_json`:

```json
{
  "api_key": "...",
  "unlocker_zone": "web_unlocker1",
  "browser_zone": "mcp_browser"
}
```

## Run

```bash
go run ./cmd/api
# GET  http://localhost:8011/health
# POST http://localhost:8011/mcp
```

```bash
docker build -t bright-data-mcp .
docker run --rm -p 8011:8011 --env-file .env bright-data-mcp
```

## Hosted Bright Data MCP (optional)

Bright Data also hosts a remote MCP at `https://mcp.brightdata.com/sse?token=…` (and `/mcp`). That path is an alternative client endpoint; **Weaver Tool Registry uses this Go HTTP server** (`mcpServer=bright-data`) with vault credentials, not the hosted URL.

## Tests

```bash
go test ./...
```

## Notes

- `web_data_*` tools call `POST /datasets/v3/trigger` and poll snapshot (same as official MCP).
- `extract` returns markdown; use your agent LLM to structure JSON (official MCP uses MCP sampling).
- Browser tools require `BRIGHTDATA_BROWSER_ZONE` and Pro mode or `GROUPS=browser`.
- DatumBridge extras still available in Pro via datasets: `web_data_reddit_comments`, `web_data_reuter_news` (not in official group lists).

## Non-goals

Native proxy access, Dataset Marketplace browsing UI, Scraper Studio IDE, replacing this server with hosted `mcp.brightdata.com`.
