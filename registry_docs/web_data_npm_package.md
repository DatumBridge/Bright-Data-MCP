# web_data_npm_package

Quickly read structured npm package data.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `package_name` | yes | Quickly read structured npm package data. |

## Cases

### Typical call

Input:

```json
{
  "package_name": "example"
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

### Missing `package_name`

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
    "error_message": "package_name is required",
    "retryable": false
  }
}
```
