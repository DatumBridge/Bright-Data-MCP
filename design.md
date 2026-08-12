# Design — bright-data-mcp

## Class

**Tool-server** (Streamable HTTP MCP), not a WS hub relay.

## Language

Go (platform exception vs default Python FastMCP), matching `shopify-mcp`.

## Upstream

Only `POST https://api.brightdata.com/request` (Web Unlocker + SERP zones).

## Auth

- Primary: vault `credentials_json` / `credentials_path`
- Fallback: `BRIGHTDATA_API_KEY` + zone env vars (local/dev)
- Zones never accepted as free-form tool arguments

## Tools

`scrape_url`, `search_serp`, `bright_data_health`

## Non-goals

Browser API, datasets, crawl orchestration, native superproxy, screenshots as first-class tools.
