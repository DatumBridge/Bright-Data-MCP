# web_data_linkedin_people_search

Quickly read structured linkedin people search data.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `url` | yes | Quickly read structured linkedin people search data. |
| `first_name` | yes | Quickly read structured linkedin people search data. |
| `last_name` | yes | Quickly read structured linkedin people search data. |

## Cases

### Typical call

Input:

```json
{
  "url": "https://example.com/page",
  "first_name": "example",
  "last_name": "example"
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

### Missing `url`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "first_name": "example",
  "last_name": "example"
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
