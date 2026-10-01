# System Overview

## Role

`bright-data-mcp` is a **stateless Go MCP tool-server** that exposes Bright Data Web Unlocker, SERP, Datasets (`web_data_*`), and Scraping Browser capabilities to DatumBridge agents via Streamable HTTP (`POST /mcp`).

Catalog parity target: [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) (Rapid / Pro / Groups).

## Layout

| Path | Responsibility |
|------|----------------|
| `cmd/api` | Process bootstrap, `/health`, CORS, logging |
| `internal/mcp` | Session store + JSON-RPC Streamable HTTP |
| `internal/brightdata` | Credentials + Direct API / datasets / browser client |
| `internal/tools` | Tool registry, Rapid/Pro/Groups filtering, handlers |

## 2026-10-01 change

- Unlocker scrape tools return the decoded Bright Data body unchanged (aside from the existing character cap) and attach the zone plus `x-brd-*` / `x-luminati-*` headers. This makes an empty Production body comparable with a Lab body for the same URL.

## 2026-09-05 change

- Removed deprecated `discover` (upstream HTTP 410).
- Rapid defaults include batch search/scrape; `session_stats` is Pro-only.
- Group membership aligned with official docs (`business` / `travel` / `research`).

## Impacted components

- MCP under `mcp/bright-data-mcp`
- Tool Registry entry `mcpServer=bright-data` (republish after deploy so `tools/list` drops `discover` and picks up per-tool `_meta.capabilities`)

## Risks

- Paid per-request Bright Data usage
- Scraped content may contain PII / prompt-injection text
- Operator ToS / legal compliance for scrape targets
- Stale Weaver Tool Registry catalogs until republish after tool-list changes

## Non-goals

- Replacing this server with hosted `https://mcp.brightdata.com/sse`
- Re-implementing Discover API
