# ADR-0001 — Direct API Bearer + Configured Zones

## Context

Bright Data supports Direct API, native proxy, Browser API, and Datasets. Agents need a scrape/search surface that stays secure-by-default (zones not injectable via tool args).

## Decision

1. Use Bright Data Direct API with Bearer API key:
   - Unlocker/SERP: `POST https://api.brightdata.com/request`
   - Datasets (`web_data_*`): trigger + snapshot poll
   - Scraping Browser: CDP via configured `browser_zone`
2. Resolve unlocker/serp/browser **zone names from credentials or env**, never free-form tool args.
3. Prefer vault `credentials_json` in multi-tenant; env keys for local/dev.
4. Do **not** call `/discover` — Discover API is retired (HTTP 410); agents use `search_engine`.
5. Filter tools with Rapid / Pro / Groups (`BRIGHTDATA_PRO_MODE`, `BRIGHTDATA_GROUPS`, `BRIGHTDATA_TOOLS`).

## Alternatives Considered

- Native superproxy — more complex TLS/cert handling
- Hosted `mcp.brightdata.com` as Weaver path — rejected; Weaver keeps self-hosted Go MCP (`mcpServer=bright-data`)
- Python FastMCP — rejected; Go + shopify-like layout required

## Consequences

Clear billing surface, zone misconfig becomes a config error not an agent injection vector, catalog stays aligned with official Bright Data MCP tools.

## Trade-offs

`extract` returns markdown for the agent LLM (official MCP sampling not mirrored). Hosted SSE MCP remains optional for external clients only.

## Risks

API shape drift; cost; ToS — mitigated by docs, truncation, and Tool Registry republish after tool-list changes.
