# Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `BRIGHTDATA_MCP_PORT` / `PORT` | `8011` | Listen port |
| `LOG_LEVEL` | `INFO` | zerolog level |
| `BRIGHTDATA_API_KEY` | — | Bearer token (env fallback) |
| `BRIGHTDATA_UNLOCKER_ZONE` | `web_unlocker2` | Web Unlocker zone (`search_engine`, scrape). Lab Direct API zone; override for other accounts. |
| `BRIGHTDATA_SERP_ZONE` | `serp_api2` | SERP zone (legacy `search_serp`). Lab Direct API zone; override for other accounts. |
| `BRIGHTDATA_BROWSER_ZONE` | — | Scraping Browser zone (Pro browser tools) |
| `BRIGHTDATA_PRO_MODE` / `PRO_MODE` | `false` | When `true`, register full tool catalog |
| `BRIGHTDATA_GROUPS` / `GROUPS` | — | CSV group ids (`ecommerce`, `social`, `browser`, …) |
| `BRIGHTDATA_TOOLS` / `TOOLS` | — | CSV explicit tool allowlist |
| `BRIGHTDATA_API_URL` | `https://api.brightdata.com/request` | Override for tests |
| `BRIGHTDATA_HTTP_TIMEOUT_SEC` | `90` | Upstream timeout |
| `BRIGHTDATA_POLLING_TIMEOUT_SEC` | — | Dataset snapshot poll budget |
| `BRIGHTDATA_MAX_CHARS` | `100000` | Default scrape truncate |
| `BRIGHTDATA_CREDENTIALS_DIR` | `/credentials` | Jail for `credentials_path` |
| `BRIGHTDATA_MCP_ALLOWED_ORIGINS` | empty | CORS allowlist (empty = no CORS) |

## Rapid vs Pro

- Rapid (default): `search_engine`, `scrape_as_markdown`, `search_engine_batch`, `scrape_batch`, plus DatumBridge aliases/health.
- Pro: all official tools including `web_data_*`, browser, `extract`, `session_stats`.
- `discover` is never registered (API removed by Bright Data).
