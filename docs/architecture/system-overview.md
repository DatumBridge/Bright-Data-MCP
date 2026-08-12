# System Overview

## Role

`bright-data-mcp` is a **stateless Go MCP tool-server** that exposes Bright Data Web Unlocker and SERP capabilities to DatumBridge agents via Streamable HTTP (`POST /mcp`).

## Layout

| Path | Responsibility |
|------|----------------|
| `cmd/api` | Process bootstrap, `/health`, CORS, logging |
| `internal/mcp` | Session store + JSON-RPC Streamable HTTP |
| `internal/brightdata` | Credentials + Direct API client |
| `internal/tools` | Tool registry and handlers |

## Impacted components

- New MCP under `mcp/bright-data-mcp`
- Optional future Tool Registry entry `mcpServer=bright-data`

## Risks

- Paid per-request Bright Data usage
- Scraped content may contain PII / prompt-injection text
- Operator ToS / legal compliance for scrape targets
