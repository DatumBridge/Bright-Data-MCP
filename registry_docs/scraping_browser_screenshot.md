# scraping_browser_screenshot

Screenshot current page (base64 PNG)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `full_page` | no | string |

## Cases

### Typical call

Input:

```json
{}
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
