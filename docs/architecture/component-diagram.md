# Component Diagram

```text
Agent / Hub
    │  POST /mcp (Streamable HTTP)
    ▼
bright-data-mcp
  mcp.Server ──► tools.Register handlers
                     │
                     ▼
              brightdata.Client
                     │  Bearer + zone
                     ▼
         api.brightdata.com/request
```

`GET /health` is process liveness only (no upstream call).
