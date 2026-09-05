# Design — bright-data-mcp

## Class

**Tool-server** (Streamable HTTP MCP), not a WS hub relay.

## Language

Go (platform exception vs default Python FastMCP), matching `shopify-mcp`.

## Upstream

- `POST https://api.brightdata.com/request` — Web Unlocker / SERP
- Datasets trigger + snapshot poll — `web_data_*`
- Scraping Browser CDP — `scraping_browser_*`
- **Not used:** `/discover` (retired; HTTP 410)

## Auth

- Primary: vault `credentials_json` / `credentials_path`
- Fallback: `BRIGHTDATA_API_KEY` + zone env vars (local/dev)
- Zones never accepted as free-form tool arguments

## Tools

Rapid/Pro/Groups catalog per [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools). DatumBridge aliases: `scrape_url`, `search_serp`, `bright_data_health`.

## Non-goals

Native superproxy, Dataset Marketplace browsing UI, Scraper Studio IDE, hosted `mcp.brightdata.com` as the Weaver integration path, re-implementing Discover.
