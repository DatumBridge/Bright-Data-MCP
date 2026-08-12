# Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `BRIGHTDATA_MCP_PORT` / `PORT` | `8011` | Listen port |
| `LOG_LEVEL` | `INFO` | zerolog level |
| `BRIGHTDATA_API_KEY` | — | Bearer token (env fallback) |
| `BRIGHTDATA_UNLOCKER_ZONE` | — | Web Unlocker zone name |
| `BRIGHTDATA_SERP_ZONE` | — | SERP zone name |
| `BRIGHTDATA_API_URL` | `https://api.brightdata.com/request` | Override for tests |
| `BRIGHTDATA_HTTP_TIMEOUT_SEC` | `90` | Upstream timeout |
| `BRIGHTDATA_MAX_CHARS` | `100000` | Default scrape truncate |
| `BRIGHTDATA_CREDENTIALS_DIR` | `/credentials` | Jail for `credentials_path` |
| `BRIGHTDATA_MCP_ALLOWED_ORIGINS` | empty | CORS allowlist (empty = no CORS) |
