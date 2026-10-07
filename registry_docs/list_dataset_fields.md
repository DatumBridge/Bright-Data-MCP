# list_dataset_fields

List filterable fields of a searchable Bright Data dataset. Supported dataset_id:

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `dataset_id` | yes | string |

## Cases

### Typical call

Input:

```json
{
  "dataset_id": "example-id"
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

### Missing `dataset_id`

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
    "error_message": "dataset_id is required",
    "retryable": false
  }
}
```
