# Integrations

## Bright Data Direct API (Unlocker / SERP)

- Endpoint: `POST https://api.brightdata.com/request`
- Auth: `Authorization: Bearer <api_key>`
- Body: `{ zone, url, format, data_format?, country?, method? }`
- Response: decoded body is returned as scrape content. Diagnostic headers (`x-brd-*`, `x-luminati-*`, content type and length) are copied into the tool metadata so an empty HTTP 200 can be compared across environments. The API key is not included.

## Bright Data Datasets API (`web_data_*`)

- Trigger: `POST https://api.brightdata.com/datasets/v3/trigger`
- Snapshot poll until ready (same pattern as official `@brightdata/mcp`)
- Catalog: `internal/brightdata/data/datasets.json`

## Discover API (removed)

- Former MCP tool `discover` called `https://api.brightdata.com/discover`
- Upstream now returns `HTTP 410: Discover API is no longer available`
- Replacement: `search_engine` / `search_engine_batch`

## Bright Data Scraping Browser

- CDP via configured `browser_zone` for `scraping_browser_*` tools

## Hosted Bright Data MCP (optional, not Weaver path)

- `https://mcp.brightdata.com/sse?token=…` / `/mcp`
- Query params: `pro=1`, `groups=…` per Bright Data docs
- DatumBridge Weaver integrates **this Go MCP**, not the hosted URL

## Related DatumBridge MCPs

- `searxng-web-search-mcp` — free/self-hosted metasearch (no Unlocker)
- `shopify-mcp` — structural Go MCP twin (not related to scraping)

## Database

N/A — no persistence.
