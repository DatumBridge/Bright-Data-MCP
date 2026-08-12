package brightdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Credentials is the vault/injected credentials_json shape for Bright Data Direct API.
// Zones are never accepted as free-form tool parameters — only from credentials/env.
type Credentials struct {
	Type         string `json:"type"`
	APIKey       string `json:"api_key"`
	UnlockerZone string `json:"unlocker_zone"`
	SerpZone     string `json:"serp_zone"`
	BrowserZone  string `json:"browser_zone"`
	// Aliases
	Token string `json:"token"`
	Key   string `json:"key"`
}

// ParseCredentials loads Credentials from vault args, then falls back to process env.
// Precedence: credentials_json / credentials_path (primary) → BRIGHTDATA_* env (local/dev).
func ParseCredentials(credentialsJSON, credentialsPath string) (*Credentials, error) {
	c := &Credentials{}
	raw := strings.TrimSpace(credentialsJSON)

	if raw == "" && strings.TrimSpace(credentialsPath) != "" {
		path := filepath.Clean(credentialsPath)
		if strings.Contains(path, "..") {
			return nil, fmt.Errorf("credentials_path must not contain ..")
		}
		allowedRoot := strings.TrimSpace(os.Getenv("BRIGHTDATA_CREDENTIALS_DIR"))
		if allowedRoot == "" {
			allowedRoot = "/credentials"
		}
		allowedRoot = filepath.Clean(allowedRoot)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		absRoot, err := filepath.Abs(allowedRoot)
		if err != nil {
			return nil, fmt.Errorf("BRIGHTDATA_CREDENTIALS_DIR: %w", err)
		}
		// Resolve symlinks so jail cannot be escaped via link outside root.
		resolvedRoot, err := filepath.EvalSymlinks(absRoot)
		if err != nil {
			// Root may not exist yet in some test setups — fall back to absRoot.
			resolvedRoot = absRoot
		}
		resolvedPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		if resolvedPath != resolvedRoot && !strings.HasPrefix(resolvedPath, resolvedRoot+string(os.PathSeparator)) {
			return nil, fmt.Errorf("credentials_path must be under %s", absRoot)
		}
		b, err := os.ReadFile(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("read credentials_path: %w", err)
		}
		raw = strings.TrimSpace(string(b))
	}

	if raw != "" {
		if err := json.Unmarshal([]byte(raw), c); err != nil {
			return nil, fmt.Errorf("invalid credentials_json: %w", err)
		}
	}

	// Env fallback / fill blanks (local/dev)
	if strings.TrimSpace(c.APIKey) == "" {
		c.APIKey = firstNonEmpty(c.Token, c.Key, os.Getenv("BRIGHTDATA_API_KEY"))
	}
	if strings.TrimSpace(c.UnlockerZone) == "" {
		c.UnlockerZone = os.Getenv("BRIGHTDATA_UNLOCKER_ZONE")
	}
	if strings.TrimSpace(c.SerpZone) == "" {
		c.SerpZone = os.Getenv("BRIGHTDATA_SERP_ZONE")
	}
	if strings.TrimSpace(c.BrowserZone) == "" {
		c.BrowserZone = firstNonEmpty(os.Getenv("BRIGHTDATA_BROWSER_ZONE"), os.Getenv("BROWSER_ZONE"))
	}

	c.APIKey = strings.TrimSpace(c.APIKey)
	c.UnlockerZone = strings.TrimSpace(c.UnlockerZone)
	c.SerpZone = strings.TrimSpace(c.SerpZone)
	c.Token = ""
	c.Key = ""

	if c.APIKey == "" {
		return nil, fmt.Errorf("api_key required (credentials_json.api_key or BRIGHTDATA_API_KEY)")
	}
	return c, nil
}

// MaskAPIKey returns a redacted preview for health responses.
func MaskAPIKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return "(missing)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "…" + key[len(key)-4:]
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
