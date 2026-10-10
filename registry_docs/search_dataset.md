# search_dataset

Search a Bright Data dataset by filter (Elasticsearch-backed). Use list_dataset_fields first.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `dataset_id` | yes | string |
| `filter` | yes | string |
| `size` | no | string |
| `sort` | no | string |

## Cases

### Typical call

Input:

```json
{
  "dataset_id": "example-id",
  "filter": "example"
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

### Missing `dataset_id`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "filter": "example"
}
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
