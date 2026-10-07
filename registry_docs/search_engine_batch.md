# search_engine_batch

Run up to 10 search queries in one request and return Markdown. queries is an array of strings, or objects with a query field. JSON organic hits are kept when Bright Data already sent JSON.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `queries` | yes | string |
| `query` | no | string |
| `engine` | no | string |
| `cursor` | no | string |
| `geo_location` | no | string |

## Cases

### Typical call

Input:

```json
{
  "queries": "example",
  "query": "quarterly report"
}
```

Output:

```json
{
  "content": [
    {
      "type": "text",
      "text": "[UNTRUSTED_WEB_CONTENT] scraped or search text"
    }
  ],
  "isError": false
}
```

### Missing `queries`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "query": "quarterly report"
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "queries is required",
    "retryable": false
  }
}
```
