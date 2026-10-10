# scrape_url

[Legacy alias] Use scrape_as_html or scrape_as_markdown. Bright Data Web Unlocker scrape.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `url` | yes | string |
| `format` | no | string |
| `data_format` | no | string |
| `country` | no | string |
| `max_chars` | no | string |

## Cases

### Typical call

Input:

```json
{
  "url": "https://example.com/page"
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
{}
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
