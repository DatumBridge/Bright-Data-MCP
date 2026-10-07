# scraping_browser_click_ref

Click element by ref from snapshot

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `ref` | yes | string |
| `element` | yes | string |

## Cases

### Typical call

Input:

```json
{
  "ref": "example",
  "element": "example"
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

### Missing `ref`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "element": "example"
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "ref is required",
    "retryable": false
  }
}
```
