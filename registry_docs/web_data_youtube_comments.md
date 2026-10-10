# web_data_youtube_comments

Quickly read structured youtube comments data.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `url` | yes | Quickly read structured youtube comments data. |
| `num_of_comments` | yes | Quickly read structured youtube comments data. |

## Cases

### Typical call

Input:

```json
{
  "url": "https://example.com/page",
  "num_of_comments": "example"
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
  "num_of_comments": "example"
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
