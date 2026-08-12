package brightdata

import (
	"bytes"
	"context"
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

const DefaultAPIURL = "https://api.brightdata.com/request"

// RequestOpts controls a Direct API /request call.
type RequestOpts struct {
	Zone       string
	URL        string
	Format     string // "raw" or "json"
	DataFormat string // "markdown" or "screenshot" when Format is "raw"
	Country    string
	Method     string
}

// Client calls Bright Data Direct API (Bearer + zone).
type Client struct {
	creds      *Credentials
	httpClient *http.Client
	apiURL     string
}

func NewClient(creds *Credentials) *Client {
	timeoutSec := 90
	if v := strings.TrimSpace(os.Getenv("BRIGHTDATA_HTTP_TIMEOUT_SEC")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutSec = n
		}
	}
	poll := pollingTimeoutSec()
	if poll+15 > timeoutSec {
		timeoutSec = poll + 15
	}
	apiURL := strings.TrimSpace(os.Getenv("BRIGHTDATA_API_URL"))
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	return &Client{
		creds: creds,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
		apiURL: apiURL,
	}
}

// Request posts to Bright Data /request and returns the response body bytes.
func (c *Client) Request(ctx context.Context, opts RequestOpts) ([]byte, int, error) {
	if c == nil || c.creds == nil {
		return nil, 0, fmt.Errorf("brightdata client not configured")
	}
	if strings.TrimSpace(opts.Zone) == "" {
		return nil, 0, fmt.Errorf("zone is required")
	}
	if err := ValidateHTTPURL(opts.URL); err != nil {
		return nil, 0, err
	}
	format := opts.Format
	if format == "" {
		format = "raw"
	}
	body := map[string]interface{}{
		"zone":   opts.Zone,
		"url":    opts.URL,
		"format": format,
	}
	if opts.DataFormat != "" {
		body["data_format"] = opts.DataFormat
	}
	if opts.Country != "" {
		body["country"] = strings.ToLower(opts.Country)
	}
	if opts.Method != "" {
		body["method"] = strings.ToUpper(opts.Method)
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.creds.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("brightdata request failed: %w", err)
	}
	defer resp.Body.Close()

	// Cap read to 30 MiB to avoid memory blowups before tool truncation.
	limited := io.LimitReader(resp.Body, 30<<20)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := sanitizeProviderError(string(data))
		return data, resp.StatusCode, fmt.Errorf("brightdata HTTP %d: %s", resp.StatusCode, msg)
	}
	return data, resp.StatusCode, nil
}

// UnlockerZone returns the configured Web Unlocker zone.
func (c *Client) UnlockerZone() string {
	if c == nil || c.creds == nil {
		return ""
	}
	return c.creds.UnlockerZone
}

// BrowserZone returns configured Scraping Browser zone.
func (c *Client) BrowserZone() string {
	if c == nil || c.creds == nil {
		return ""
	}
	return c.creds.BrowserZone
}

// SerpZone returns the configured SERP zone.
func (c *Client) SerpZone() string {
	if c == nil || c.creds == nil {
		return ""
	}
	return c.creds.SerpZone
}

// APIKeyMasked returns a redacted API key preview.
func (c *Client) APIKeyMasked() string {
	if c == nil || c.creds == nil {
		return "(missing)"
	}
	return MaskAPIKey(c.creds.APIKey)
}

// ValidateHTTPURL ensures url is http(s) with a host.
func ValidateHTTPURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("url is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("url must use http or https scheme")
	}
	if u.Host == "" {
		return fmt.Errorf("url must include a host")
	}
	return nil
}

// BuildSERPURL builds a Google/Bing search URL for the SERP zone.
func BuildSERPURL(query, engine, country, language string, start int, brdJSON bool) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	if len(query) > 500 {
		return "", fmt.Errorf("query exceeds 500 characters")
	}
	engine = strings.ToLower(strings.TrimSpace(engine))
	if engine == "" {
		engine = "google"
	}
	country = strings.ToLower(strings.TrimSpace(country))
	if country == "" {
		country = "us"
	}
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		language = "en"
	}
	if start < 0 {
		start = 0
	}
	if start > 1000 {
		start = 1000
	}

	switch engine {
	case "google":
		u, _ := url.Parse("https://www.google.com/search")
		q := u.Query()
		q.Set("q", query)
		q.Set("hl", language)
		q.Set("gl", country)
		q.Set("start", strconv.Itoa(start))
		if brdJSON {
			q.Set("brd_json", "1")
		}
		u.RawQuery = q.Encode()
		return u.String(), nil
	case "bing":
		u, _ := url.Parse("https://www.bing.com/search")
		q := u.Query()
		q.Set("q", query)
		q.Set("setlang", language)
		q.Set("cc", country)
		q.Set("first", strconv.Itoa(start+1))
		if brdJSON {
			q.Set("brd_json", "1")
		}
		u.RawQuery = q.Encode()
		return u.String(), nil
	default:
		return "", fmt.Errorf("unsupported engine %q (use google or bing)", engine)
	}
}

// TruncateUTF8 truncates s to at most maxChars runes (approx chars).
func TruncateUTF8(s string, maxChars int) (out string, truncated bool) {
	if maxChars <= 0 {
		return "", true
	}
	runes := []rune(s)
	if len(runes) <= maxChars {
		return s, false
	}
	return string(runes[:maxChars]), true
}

func sanitizeProviderError(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return "(empty body)"
	}
	// Avoid dumping huge error HTML into agent context.
	if len(body) > 500 {
		body = body[:500] + "…"
	}
	// Best-effort: strip likely token-looking substrings.
	lower := strings.ToLower(body)
	if strings.Contains(lower, "bearer ") {
		return "(provider error redacted)"
	}
	return body
}

// DefaultMaxChars returns global default truncate size.
func DefaultMaxChars() int {
	if v := strings.TrimSpace(os.Getenv("BRIGHTDATA_MAX_CHARS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 100_000
}
