# web_data_perplexity_ai_insights

Send a prompt to Perplexity and get back AI-generated insights.

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `prompt` | yes | Send a prompt to Perplexity and get back AI-generated insights. |

## Cases

### Typical call

Input:

```json
{
  "prompt": "Summarize the main offer on this page"
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

### Missing `prompt`

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
    "error_message": "prompt is required",
    "retryable": false
  }
}
```
