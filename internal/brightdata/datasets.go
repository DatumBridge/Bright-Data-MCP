package brightdata

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed data/datasets.json
var datasetsJSON []byte

// DatasetTool describes a Bright Data web_data_* MCP tool (from official @brightdata/mcp catalog).
type DatasetTool struct {
	ID             string            `json:"id"`
	DatasetID      string            `json:"dataset_id"`
	Description    string            `json:"description"`
	Inputs         []string          `json:"inputs"`
	Defaults       map[string]string `json:"defaults,omitempty"`
	FixedValues    map[string]string `json:"fixed_values,omitempty"`
	TriggerParams  map[string]string `json:"trigger_params,omitempty"`
}

// LoadDatasetCatalog returns the embedded web_data tool catalog.
func LoadDatasetCatalog() ([]DatasetTool, error) {
	var raw []DatasetTool
	if err := json.Unmarshal(datasetsJSON, &raw); err != nil {
		return nil, err
	}
	applyDatasetOverrides(&raw)
	return raw, nil
}

func applyDatasetOverrides(tools *[]DatasetTool) {
	byID := map[string]*DatasetTool{}
	for i := range *tools {
		byID[(*tools)[i].ID] = &(*tools)[i]
	}
	set := func(id string, fn func(*DatasetTool)) {
		if t, ok := byID[id]; ok {
			fn(t)
		}
	}
	set("amazon_product_search", func(t *DatasetTool) {
		t.FixedValues = map[string]string{"pages_to_search": "1"}
	})
	set("google_maps_reviews", func(t *DatasetTool) {
		t.Defaults = map[string]string{"days_limit": "3"}
	})
	set("youtube_comments", func(t *DatasetTool) {
		t.Defaults = map[string]string{"num_of_comments": "10"}
	})
	set("reddit_comments", func(t *DatasetTool) {
		t.Defaults = map[string]string{"days_back": ""}
	})
	set("x_profile_posts", func(t *DatasetTool) {
		t.Defaults = map[string]string{"start_date": "", "end_date": ""}
		t.TriggerParams = map[string]string{
			"type":              "discover_new",
			"discover_by":       "profile_url_most_recent_posts",
			"limit_per_input":   "10",
		}
	})
	set("chatgpt_ai_insights", func(t *DatasetTool) {
		t.FixedValues = map[string]string{
			"url": "https://chatgpt.com/", "country": "", "web_search": "false", "additional_prompt": "",
		}
		t.TriggerParams = map[string]string{"custom_output_fields": "answer_text_markdown"}
	})
	set("grok_ai_insights", func(t *DatasetTool) {
		t.FixedValues = map[string]string{"url": "https://grok.com/", "index": ""}
		t.TriggerParams = map[string]string{"custom_output_fields": "answer_text_markdown"}
	})
	set("perplexity_ai_insights", func(t *DatasetTool) {
		t.FixedValues = map[string]string{"url": "https://www.perplexity.ai", "index": "", "country": ""}
		t.TriggerParams = map[string]string{"custom_output_fields": "answer_text_markdown"}
	})
}

// TriggerWebData runs datasets/v3/trigger and polls snapshot until ready.
func (c *Client) TriggerWebData(ctx context.Context, tool DatasetTool, input map[string]string) ([]byte, error) {
	if c == nil || c.creds == nil {
		return nil, fmt.Errorf("brightdata client not configured")
	}
	payload := map[string]interface{}{}
	for k, v := range input {
		if strings.TrimSpace(v) != "" {
			payload[k] = v
		}
	}
	for k, v := range tool.FixedValues {
		payload[k] = v
	}
	for k, def := range tool.Defaults {
		if _, ok := payload[k]; !ok {
			payload[k] = def
		}
	}
	// AI insight tools use prompt as primary input.
	if tool.ID == "chatgpt_ai_insights" || tool.ID == "grok_ai_insights" || tool.ID == "perplexity_ai_insights" {
		if p, ok := payload["prompt"]; ok {
			payload["prompt"] = p
		}
	}

	q := url.Values{}
	q.Set("dataset_id", tool.DatasetID)
	q.Set("include_errors", "true")
	for k, v := range tool.TriggerParams {
		q.Set(k, v)
	}
	triggerURL := "https://api.brightdata.com/datasets/v3/trigger?" + q.Encode()
	body, err := json.Marshal([]map[string]interface{}{payload})
	if err != nil {
		return nil, err
	}
	var triggerResp struct {
		SnapshotID string `json:"snapshot_id"`
	}
	if err := c.doJSON(ctx, http.MethodPost, triggerURL, body, &triggerResp); err != nil {
		return nil, fmt.Errorf("dataset trigger: %w", err)
	}
	if triggerResp.SnapshotID == "" {
		return nil, fmt.Errorf("dataset trigger: no snapshot_id returned")
	}
	return c.pollDatasetSnapshot(ctx, triggerResp.SnapshotID)
}

func (c *Client) pollDatasetSnapshot(ctx context.Context, snapshotID string) ([]byte, error) {
	maxAttempts := pollingTimeoutSec()
	snapshotURL := fmt.Sprintf("https://api.brightdata.com/datasets/v3/snapshot/%s?format=json", snapshotID)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, snapshotURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.creds.APIKey)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 30<<20))
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("snapshot poll HTTP %d: %s", resp.StatusCode, sanitizeProviderError(string(data)))
		}
		var statusWrap struct {
			Status string `json:"status"`
		}
		_ = json.Unmarshal(data, &statusWrap)
		if statusWrap.Status == "running" || statusWrap.Status == "building" || statusWrap.Status == "starting" {
			time.Sleep(time.Second)
			continue
		}
		return data, nil
	}
	return nil, fmt.Errorf("timeout after %d seconds waiting for dataset snapshot", maxAttempts)
}

// PollingTimeoutSec returns max poll attempts (1 second each) for dataset APIs.
func PollingTimeoutSec() int {
	return pollingTimeoutSec()
}

func pollingTimeoutSec() int {
	if v := strings.TrimSpace(os.Getenv("BRIGHTDATA_POLLING_TIMEOUT_SEC")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 120
}

func (c *Client) doJSON(ctx context.Context, method, rawURL string, body []byte, out interface{}) error {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.creds.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 30<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, sanitizeProviderError(string(data)))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(data, out)
}

// DatasetMetadata returns filterable fields for searchable datasets.
func (c *Client) DatasetMetadata(ctx context.Context, datasetID string) ([]byte, error) {
	u := fmt.Sprintf("https://api.brightdata.com/datasets/%s/metadata", datasetID)
	var out json.RawMessage
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchDataset queries datasets/search/{id}.
func (c *Client) SearchDataset(ctx context.Context, datasetID string, body map[string]interface{}) ([]byte, error) {
	u := fmt.Sprintf("https://api.brightdata.com/datasets/search/%s", datasetID)
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.creds.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 30<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, sanitizeProviderError(string(data)))
	}
	return data, nil
}
