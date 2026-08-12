package brightdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// BrowserCDPEndpoint resolves Scraping Browser CDP WebSocket URL (official MCP parity).
func BrowserCDPEndpoint(c *Client, zone, country string) (string, error) {
	if c == nil || c.creds == nil {
		return "", fmt.Errorf("client not configured")
	}
	var status struct {
		Customer string `json:"customer"`
	}
	if err := c.doJSON(context.Background(), http.MethodGet, "https://api.brightdata.com/status", nil, &status); err != nil {
		return "", fmt.Errorf("browser status: %w", err)
	}
	pwURL := fmt.Sprintf("https://api.brightdata.com/zone/passwords?zone=%s", zone)
	req, err := http.NewRequest(http.MethodGet, pwURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.creds.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var pwResp struct {
		Passwords []string `json:"passwords"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pwResp); err != nil {
		return "", err
	}
	if len(pwResp.Passwords) == 0 {
		return "", fmt.Errorf("no browser zone password for zone %q", zone)
	}
	countrySuffix := ""
	if country != "" {
		countrySuffix = "-country-" + country
	}
	return fmt.Sprintf("wss://brd-customer-%s-zone-%s%s:%s@brd.superproxy.io:9222",
		status.Customer, zone, countrySuffix, pwResp.Passwords[0]), nil
}
