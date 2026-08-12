# Security Architecture

## Trust model

- Bright Data API key authorizes billable Unlocker/SERP requests as the operator account
- Prefer vault `credentials_json` so keys are not process-wide in multi-tenant
- Scraped/search bodies are **untrusted** — prefixed `[UNTRUSTED_WEB_CONTENT]`

## Controls

| Control | Behavior |
|---------|----------|
| Zones from config only | Tools cannot pick arbitrary Bright Data zones |
| URL scheme allowlist | `http`/`https` only |
| Truncation | Default `max_chars=100000` (cap 500000) |
| Health | Masks API key; `probe=false` by default |
| Logging | Never log Bearer tokens or full scrape bodies at info |
| Path jail | `credentials_path` must stay under `BRIGHTDATA_CREDENTIALS_DIR` |

## Operator duties

Comply with target ToS and law. Do not scrape authenticated/private content without authorization.
