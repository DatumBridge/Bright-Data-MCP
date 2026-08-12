# ADR-0001 — Direct API Bearer + Configured Zones

## Context

Bright Data supports Direct API, native proxy, Browser API, and Datasets. Agents need a small scrape/search surface.

## Decision

1. Use only `POST https://api.brightdata.com/request` with Bearer API key.
2. Resolve unlocker/serp **zone names from credentials or env**, never free-form tool args.
3. Prefer vault `credentials_json` in multi-tenant; env keys for local/dev.
4. Defer Browser/Datasets/crawl to later ADRs.

## Alternatives Considered

- Native superproxy — more complex TLS/cert handling
- Browser API — overkill for v1 URL scrape
- Python FastMCP — rejected; user required Go + shopify-like layout

## Consequences

Simple client, clear billing surface, zone misconfig becomes a config error not an agent injection vector.

## Trade-offs

No interactive browsing or dataset library scrapers in v1.

## Risks

API shape drift; cost; ToS — mitigated by docs and truncation.
