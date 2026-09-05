# Use Cases

| Actor | Goal | Tool |
|-------|------|------|
| Agent | Search the public web | `search_engine` (or `search_engine_batch`) |
| Agent | Read a public page as markdown | `scrape_as_markdown` (or `scrape_batch`) |
| Agent | Read raw HTML | `scrape_as_html` (Pro / `advanced_scraping`) |
| Agent | Structured Amazon/LinkedIn/… data | `web_data_*` (Pro / group) |
| Agent | Browser automation | `scraping_browser_*` (Pro / `browser`) |
| Agent | Legacy scrape/search | `scrape_url` / `search_serp` |
| Operator | Verify config before demos | `bright_data_health` |

## Do not use

| Tool | Status |
|------|--------|
| `discover` | Removed — Discover API HTTP 410; use `search_engine` |
