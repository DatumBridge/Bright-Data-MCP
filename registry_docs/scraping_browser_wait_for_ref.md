# scraping_browser_wait_for_ref

Wait for ref element visible

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `ref` | yes | string |
| `timeout_ms` | no | string |

## Cases

### Typical call

Input:

```json
{
  "ref": "example"
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
{}
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
