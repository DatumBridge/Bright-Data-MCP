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

### Unlocker response

`scrape_as_markdown`, `scrape_as_html`, `scrape_url`, and `extract` return one text block:

1. A JSON object with the request zone, format, HTTP status, `char_count`, `empty_body`, and `brightdata_headers`.
2. A `--- content ---` separator.
3. The decoded Bright Data body, unchanged apart from `BRIGHTDATA_MAX_CHARS` truncation (`truncated: true` when that applies).

`brightdata_headers` includes `x-brd-*`, `x-luminati-*`, `content-type`, and `content-length`. Cookies and authorization headers are not copied. `success` is false when `x-brd-error` or `x-luminati-error` is present. `scrape_batch` items carry the same status, zone, headers, and exact `content`.

Compare `zone` and `brightdata_headers` when Lab and Production return different bodies for the same URL. The zone comes from `credentials_json.unlocker_zone`, then `BRIGHTDATA_UNLOCKER_ZONE`.

#### search_serp

Args: `query*` , `engine` (`google`\|`bing`), `country`, `language`, `start`, `brd_json` (default `false`), `max_results`, credentials…

`brd_json=false` returns the fetched page. `brd_json=true` asks Bright Data to parse Google into JSON and is slower. Prefer `search_engine`, which always requests Markdown and returns that text when the body is not JSON.

#### bright_data_health

Args: `probe` (default false), credentials…

Domain failures return MCP `isError: true`. Protocol errors use JSON-RPC codes.
