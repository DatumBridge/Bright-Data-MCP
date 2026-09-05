# API Specification

## HTTP

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Liveness JSON |
| POST | `/mcp` | Streamable HTTP MCP (protocol `2024-11-05`) |

## Tool catalog

Full names and groups: [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) and `docs/technical/MCP_TOOL_SPECIFICATION.md`.

### Modes

- **Rapid (default):** `search_engine`, `scrape_as_markdown`, `search_engine_batch`, `scrape_batch`, plus DatumBridge aliases `scrape_url` / `search_serp` / `bright_data_health`
- **Pro (`BRIGHTDATA_PRO_MODE=true`):** 60+ tools including `web_data_*`, `scraping_browser_*`, `extract`, `session_stats`, `list_dataset_fields`, `search_dataset`
- **Groups / Tools:** `BRIGHTDATA_GROUPS` / `BRIGHTDATA_TOOLS` CSV filters

### Removed

- `discover` — Bright Data Discover API returns HTTP 410; use `search_engine`

### Legacy aliases (still registered in Rapid)

#### scrape_url

Args: `url*` , `format` (`raw`\|`json`, default `raw`), `data_format` (`markdown`\|`screenshot`, only when `format=raw`), `country`, `max_chars`, credentials…

Maps 1:1 to Bright Data `POST /request` body: `{ zone, url, format, data_format?, country? }`. Prefer `scrape_as_markdown` / `scrape_as_html`.

#### search_serp

Args: `query*` , `engine` (`google`\|`bing`), `country`, `language`, `start`, `brd_json`, `max_results`, credentials…

Prefer `search_engine`.

#### bright_data_health

Args: `probe` (default false), credentials…

Domain failures return MCP `isError: true`. Protocol errors use JSON-RPC codes.
