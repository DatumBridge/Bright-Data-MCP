# Workflows

## Local connect

1. Create Unlocker (+ SERP) zones; copy API key
2. Set `.env` from `.env.example`
3. `go run ./cmd/api`
4. `initialize` → `tools/list` → `tools/call`

## Scrape page

1. Agent calls `scrape_url` with `url` (+ optional `format`, `country`, `max_chars`)
2. Server uses Unlocker zone + Bearer key
3. Returns truncated markdown/html with untrusted prefix

## SERP search

1. Agent calls `search_serp` with `query`
2. Server builds Google/Bing URL (`brd_json=1` by default)
3. Returns capped organic results

## Health

1. Call `bright_data_health` (no upstream)
2. Optionally `probe=true` for a billable Unlocker check against `https://example.com`
