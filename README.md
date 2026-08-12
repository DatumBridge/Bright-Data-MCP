# Bright Data MCP Server

Go MCP **tool-server** aligned with [Bright Data MCP tools](https://docs.brightdata.com/ai/mcp-server/tools) (60+ tools).

**mcpServer id:** `bright-data` (aliases: `bright-data-mcp`)

## Tool modes

| Mode | Env | Tools |
|------|-----|-------|
| **Rapid (default)** | `BRIGHTDATA_PRO_MODE=false` | `search_engine`, `scrape_as_markdown`, `discover`, `session_stats`, legacy `scrape_url` / `search_serp`, `bright_data_health` |
| **Pro** | `BRIGHTDATA_PRO_MODE=true` | All 60+ tools including `web_data_*`, browser automation, batch scrape, dataset search |
| **Groups** | `BRIGHTDATA_GROUPS=ecommerce,social` | Subset per [official groups](https://github.com/brightdata/brightdata-mcp) |
| **Custom** | `BRIGHTDATA_TOOLS=tool1,tool2` | Explicit allowlist |

## Tool categories (Pro)

- **Search/scrape:** `search_engine`, `scrape_as_markdown`, `scrape_as_html`, `scrape_batch`, `search_engine_batch`, `extract`, `discover`
- **Structured data:** `web_data_amazon_product`, `web_data_linkedin_person_profile`, … (50 datasets)
- **Dataset search:** `list_dataset_fields`, `search_dataset`
- **Browser:** `scraping_browser_*` (13 tools via Scraping Browser + CDP)
- **Legacy aliases:** `scrape_url`, `search_serp`

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

## Tests

```bash
go test ./...
```

## Notes

- `web_data_*` tools call `POST /datasets/v3/trigger` and poll snapshot (same as official MCP).
- `extract` returns markdown; use your agent LLM to structure JSON (official MCP uses MCP sampling).
- Browser tools require `BRIGHTDATA_BROWSER_ZONE` and Pro mode or `GROUPS=browser`.

## Non-goals

Native proxy access, Dataset Marketplace browsing UI, Scraper Studio IDE.
