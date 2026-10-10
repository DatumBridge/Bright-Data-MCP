# search_serp

[Legacy alias] Use search_engine. SERP via Web Unlocker. brd_json defaults to false so the page is returned without Bright Data JSON parsing. Set brd_json=true to request parsed JSON.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `query` | yes | string |
| `engine` | no | string |
| `country` | no | string |
| `language` | no | string |
| `start` | no | string |
| `brd_json` | no | string |
| `max_results` | no | string |

## Cases

### Typical call

Input:

```json
{
  "query": "quarterly report"
}
```

Output:

```json
{
  "content": [
    {
      "type": "text",
      "text": "scraped or search text"
    }
  ],
  "isError": false
}
```

### Missing `query`

The tool rejects the call and does not guess the missing value.

Input:

```json
{}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "query is required",
    "retryable": false
  }
}
```
