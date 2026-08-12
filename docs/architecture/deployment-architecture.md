# Deployment Architecture

- Binary: `cmd/api` or Docker image `bright-data-mcp`
- Default port: **8011**
- Probes: `GET /health` → `{"status":"ok","service":"bright-data-mcp"}`
- Secrets: inject via env (dev) or vault `credentials_json` at tool-call time
- No database; in-memory MCP sessions (30m TTL) — single replica affinity like shopify-mcp
