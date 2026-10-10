# search_engine

Scrape search results from Google, Bing, or Yandex as Markdown. JSON organic results are returned only when Bright Data already sent JSON.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `query` | yes | string |
| `engine` | no | string |
| `cursor` | no | string |
| `geo_location` | no | string |

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
