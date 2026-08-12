package brightdata_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
)

func TestParseCredentialsFromJSON(t *testing.T) {
	c, err := brightdata.ParseCredentials(`{"api_key":"secretkey123","unlocker_zone":"wu1","serp_zone":"serp1"}`, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "secretkey123" || c.UnlockerZone != "wu1" || c.SerpZone != "serp1" {
		t.Fatalf("unexpected creds: %+v", c)
	}
}

func TestParseCredentialsEnvFallback(t *testing.T) {
	t.Setenv("BRIGHTDATA_API_KEY", "envkey")
	t.Setenv("BRIGHTDATA_UNLOCKER_ZONE", "env-wu")
	t.Setenv("BRIGHTDATA_SERP_ZONE", "env-serp")
	c, err := brightdata.ParseCredentials("", "")
	if err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "envkey" || c.UnlockerZone != "env-wu" || c.SerpZone != "env-serp" {
		t.Fatalf("unexpected: %+v", c)
	}
}

func TestParseCredentialsPathJail(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BRIGHTDATA_CREDENTIALS_DIR", root)
	path := filepath.Join(root, "creds.json")
	if err := os.WriteFile(path, []byte(`{"api_key":"fromfile","unlocker_zone":"z"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := brightdata.ParseCredentials("", path)
	if err != nil {
		t.Fatal(err)
	}
	if c.APIKey != "fromfile" {
		t.Fatalf("got %q", c.APIKey)
	}

	outside := filepath.Join(t.TempDir(), "escape.json")
	_ = os.WriteFile(outside, []byte(`{"api_key":"x"}`), 0o600)
	if _, err := brightdata.ParseCredentials("", outside); err == nil {
		t.Fatal("expected jail error")
	}
}

func TestMaskAPIKey(t *testing.T) {
	if brightdata.MaskAPIKey("") != "(missing)" {
		t.Fatal("empty")
	}
	m := brightdata.MaskAPIKey("abcdefghijklmnop")
	if m == "abcdefghijklmnop" || !strings.Contains(m, "…") {
		t.Fatalf("mask failed: %s", m)
	}
}

func TestValidateHTTPURL(t *testing.T) {
	if err := brightdata.ValidateHTTPURL("https://example.com/a"); err != nil {
		t.Fatal(err)
	}
	if err := brightdata.ValidateHTTPURL("ftp://example.com"); err == nil {
		t.Fatal("expected scheme error")
	}
	if err := brightdata.ValidateHTTPURL("not-a-url"); err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildSERPURL(t *testing.T) {
	u, err := brightdata.BuildSERPURL("hello world", "google", "us", "en", 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(u, "google.com") || !contains(u, "brd_json=1") || !contains(u, "q=hello") {
		t.Fatalf("unexpected url: %s", u)
	}
	u2, err := brightdata.BuildSERPURL("x", "bing", "", "", 10, false)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(u2, "bing.com") {
		t.Fatalf("bing: %s", u2)
	}
	if _, err := brightdata.BuildSERPURL("", "google", "", "", 0, true); err == nil {
		t.Fatal("empty query")
	}
}

func TestTruncateUTF8(t *testing.T) {
	out, trunc := brightdata.TruncateUTF8("abcdef", 3)
	if out != "abc" || !trunc {
		t.Fatalf("%q %v", out, trunc)
	}
	out2, trunc2 := brightdata.TruncateUTF8("ab", 10)
	if out2 != "ab" || trunc2 {
		t.Fatalf("%q %v", out2, trunc2)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
