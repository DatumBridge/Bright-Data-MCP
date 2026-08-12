package brightdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// DiscoverRequest triggers Bright Data /discover.
type DiscoverRequest struct {
	Query            string   `json:"query"`
	Format           string   `json:"format"`
	Intent           string   `json:"intent,omitempty"`
	Country          string   `json:"country,omitempty"`
	City             string   `json:"city,omitempty"`
	Language         string   `json:"language,omitempty"`
	NumResults       int      `json:"num_results,omitempty"`
	FilterKeywords   []string `json:"filter_keywords,omitempty"`
	RemoveDuplicates *bool    `json:"remove_duplicates,omitempty"`
	StartDate        string   `json:"start_date,omitempty"`
	EndDate          string   `json:"end_date,omitempty"`
}

// Discover polls /discover until results are ready.
func (c *Client) Discover(ctx context.Context, req DiscoverRequest) ([]byte, error) {
	if req.Format == "" {
		req.Format = "json"
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var trigger struct {
		TaskID string `json:"task_id"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "https://api.brightdata.com/discover", body, &trigger); err != nil {
		return nil, err
	}
	if trigger.TaskID == "" {
		return nil, fmt.Errorf("discover: no task_id returned")
	}
	maxAttempts := pollingTimeoutSec()
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		u := "https://api.brightdata.com/discover?" + url.Values{"task_id": {trigger.TaskID}}.Encode()
		var poll struct {
			Status  string `json:"status"`
			Results []struct {
				Link            string  `json:"link"`
				Title           string  `json:"title"`
				Description     string  `json:"description"`
				RelevanceScore  float64 `json:"relevance_score"`
			} `json:"results"`
		}
		if err := c.doJSON(ctx, http.MethodGet, u, nil, &poll); err != nil {
			if attempt+1 >= maxAttempts {
				return nil, err
			}
			time.Sleep(time.Second)
			continue
		}
		if poll.Status == "processing" {
			time.Sleep(time.Second)
			continue
		}
		out := make([]map[string]interface{}, 0, len(poll.Results))
		for _, r := range poll.Results {
			out = append(out, map[string]interface{}{
				"link":             r.Link,
				"title":            r.Title,
				"description":      r.Description,
				"relevance_score":  r.RelevanceScore,
			})
		}
		return json.Marshal(out)
	}
	return nil, fmt.Errorf("discover timeout after %d seconds", maxAttempts)
}
