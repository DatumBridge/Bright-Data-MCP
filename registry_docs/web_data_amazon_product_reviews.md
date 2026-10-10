# web_data_amazon_product_reviews

Quickly read structured amazon product review data. Requires a valid product URL with /dp/ in it.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `url` | yes | Quickly read structured amazon product review data. Requires a valid product URL with /dp/ in it. |

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
