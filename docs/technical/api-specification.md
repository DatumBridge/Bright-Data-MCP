# API Specification

## HTTP

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Liveness JSON |
| POST | `/mcp` | Streamable HTTP MCP (protocol `2024-11-05`) |

## Tools

### scrape_url

Args: `url*` , `format` (`raw`\|`json`, default `raw`), `data_format` (`markdown`\|`screenshot`, only when `format=raw`), `country`, `max_chars`, credentials…

Maps 1:1 to Bright Data `POST /request` body: `{ zone, url, format, data_format?, country? }`.

### search_serp

Args: `query*` , `engine` (`google`\|`bing`), `country`, `language`, `start`, `brd_json`, `max_results`, credentials…

### bright_data_health

Args: `probe` (default false), credentials…

Domain failures return MCP `isError: true`. Protocol errors use JSON-RPC codes.
