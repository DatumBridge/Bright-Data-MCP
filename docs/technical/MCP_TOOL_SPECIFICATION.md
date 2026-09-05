# MCP Tool Specification

| Field | Value |
|-------|-------|
| Folder | `bright-data-mcp` |
| mcpServer id | `bright-data` |
| Aliases | `bright-data-mcp` |
| Transport | Streamable HTTP `POST /mcp` |
| Protocol | `2024-11-05` |
| Health | `GET /health` |
| Default port | `8011` |
| Catalog | [Bright Data MCP tools](https://docs.brightdata.com/products/mcp-server/tools) |

## Modes

| Mode | Env | Exposed tools |
|------|-----|---------------|
| Rapid (default) | `BRIGHTDATA_PRO_MODE` unset/false | `search_engine`, `scrape_as_markdown`, `search_engine_batch`, `scrape_batch`, `scrape_url`, `search_serp`, `bright_data_health` |
| Pro | `BRIGHTDATA_PRO_MODE=true` | Full catalog (60+), including `web_data_*`, browser, `extract`, `session_stats`, dataset search |
| Groups | `BRIGHTDATA_GROUPS=…` | Official group subsets (`ecommerce`, `social`, `browser`, …) |
| Custom | `BRIGHTDATA_TOOLS=…` | Explicit allowlist |

## Rapid tool names

- `search_engine`
- `scrape_as_markdown`
- `search_engine_batch`
- `scrape_batch`
- `scrape_url` (legacy alias)
- `search_serp` (legacy alias)
- `bright_data_health`

## Pro highlights

- Advanced scrape: `scrape_as_html`, `extract`, `session_stats`
- Dataset search: `list_dataset_fields`, `search_dataset`
- Structured: `web_data_*` from `internal/brightdata/data/datasets.json`
- Browser: `scraping_browser_*` (13)

## Removed

| Tool | Reason | Replacement |
|------|--------|-------------|
| `discover` | Bright Data Discover API `HTTP 410` | `search_engine` / `search_engine_batch` |

## Hosted MCP (out of band)

`https://mcp.brightdata.com/sse?token=…` and `/mcp` are Bright Data–hosted alternatives. DatumBridge Weaver integrates this Go server, not the hosted URL.
