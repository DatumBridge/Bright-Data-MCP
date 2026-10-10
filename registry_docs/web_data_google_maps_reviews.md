# web_data_google_maps_reviews

Quickly read structured Google maps reviews data.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `url` | yes | Quickly read structured Google maps reviews data. |
| `days_limit` | yes | Quickly read structured Google maps reviews data. |

## Cases

### Typical call

Input:

```json
{
  "url": "https://example.com/page",
  "days_limit": "example"
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

### Missing `url`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "days_limit": "example"
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "url is required",
    "retryable": false
  }
}
```
