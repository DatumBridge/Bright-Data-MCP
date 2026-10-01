package brightdata_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
)

func TestClientRequestSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("auth header %q", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"zone":"wu1"`) {
			t.Errorf("body %s", body)
		}
		w.Header().Set("X-Brd-Error", "render_failed")
		w.Header().Set("X-Brd-Status-Code", "200")
		w.Header().Set("Set-Cookie", "session=secret")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# Hello"))
	}))
	defer srv.Close()

	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	c := brightdata.NewClient(&brightdata.Credentials{APIKey: "test-key", UnlockerZone: "wu1"})
	resp, err := c.Request(context.Background(), brightdata.RequestOpts{
		Zone:       "wu1",
		URL:        "https://example.com",
		Format:     "raw",
		DataFormat: "markdown",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != 200 || string(resp.Body) != "# Hello" {
		t.Fatalf("status=%d data=%q", resp.Status, resp.Body)
	}
	if got := brightdata.DiagnosticHeaders(resp.Header)["x-brd-error"]; got != "render_failed" {
		t.Fatalf("headers %#v", got)
	}
	if _, ok := brightdata.DiagnosticHeaders(resp.Header)["set-cookie"]; ok {
		t.Fatal("credential-bearing header leaked")
	}
}

func TestClientRequestHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"error":"forbidden"}`))
	}))
	defer srv.Close()
	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	c := brightdata.NewClient(&brightdata.Credentials{APIKey: "k"})
	_, err := c.Request(context.Background(), brightdata.RequestOpts{
		Zone: "z",
		URL:  "https://example.com",
	})
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("expected 403 error, got %v", err)
	}
}
