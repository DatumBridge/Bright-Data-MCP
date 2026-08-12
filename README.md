# Bright Data MCP Server

Go MCP **tool-server** for Bright Data Web Unlocker + SERP API. Structure mirrors `shopify-mcp`.

**mcpServer id:** `bright-data` (aliases: `bright-data-mcp`)

Complements `searxng-web-search-mcp` (metasearch). Use this server when you need Bright Data unlocking / paid SERP.

## Tools (3)

| Tool | Description |
|------|-------------|
| `scrape_url` | Scrape a public URL via Web Unlocker (`markdown` default or `html`) |
| `search_serp` | Google/Bing SERP via SERP zone (capped organic results) |
| `bright_data_health` | Masked credential/zone check; optional billable `probe=true` |

## Setup

1. Create Web Unlocker (+ optional SERP) zones in [Bright Data Control Panel](https://brightdata.com/cp)
2. Copy `.env.example` → `.env` and set:

```bash
BRIGHTDATA_API_KEY=...
BRIGHTDATA_UNLOCKER_ZONE=...
BRIGHTDATA_SERP_ZONE=...   # required for search_serp
```

Vault inject (preferred in multi-tenant):

```json
{
  "type": "api_key",
  "api_key": "...",
  "unlocker_zone": "web_unlocker1",
  "serp_zone": "serp_api1"
}
```

Pass as `credentials_json` / `credentials_path` on each tool call.

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

## Compliance

Operators must comply with target site Terms of Service and applicable law. Scraped/search content is prefixed as untrusted for agents. Prefer `format=markdown` and `max_chars` to limit model context.

## Non-goals (v1)

Browser API, Dataset marketplace, crawl jobs, native proxy access, screenshots.
