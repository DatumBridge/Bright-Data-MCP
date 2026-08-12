# API Specification

## HTTP

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Liveness JSON |
| POST | `/mcp` | Streamable HTTP MCP (protocol `2024-11-05`) |

## Tools

### scrape_url

Args: `url*` , `format` (`markdown`\|`html`), `country`, `max_chars`, `credentials_json`, `credentials_path`

### search_serp

Args: `query*` , `engine` (`google`\|`bing`), `country`, `language`, `start`, `brd_json`, `max_results`, credentials…

### bright_data_health

Args: `probe` (default false), credentials…

Domain failures return MCP `isError: true`. Protocol errors use JSON-RPC codes.
