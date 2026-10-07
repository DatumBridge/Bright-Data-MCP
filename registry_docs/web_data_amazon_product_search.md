# web_data_amazon_product_search

Quickly read structured amazon product search data. Requires a search keyword and amazon domain URL.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `keyword` | yes | Quickly read structured amazon product search data. Requires a search keyword and amazon domain URL. |
| `url` | yes | Quickly read structured amazon product search data. Requires a search keyword and amazon domain URL. |

## Cases

### Typical call

Input:

```json
{
  "keyword": "running shoes",
  "url": "https://example.com/page"
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

### Missing `keyword`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "url": "https://example.com/page"
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "keyword is required",
    "retryable": false
  }
}
```
