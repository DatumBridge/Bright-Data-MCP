# scrape_batch

Scrape up to 10 webpages and return URL/markdown pairs.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `urls` | yes | string |

## Cases

### Typical call

Input:

```json
{
  "urls": [
    "https://example.com/a",
    "https://example.com/b"
  ]
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

### Missing `urls`

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
    "error_message": "urls is required",
    "retryable": false
  }
}
```
