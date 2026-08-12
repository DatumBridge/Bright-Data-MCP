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
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# Hello"))
	}))
	defer srv.Close()

	t.Setenv("BRIGHTDATA_API_URL", srv.URL)
	c := brightdata.NewClient(&brightdata.Credentials{APIKey: "test-key", UnlockerZone: "wu1"})
	data, status, err := c.Request(context.Background(), brightdata.RequestOpts{
		Zone:       "wu1",
		URL:        "https://example.com",
		Format:     "raw",
		DataFormat: "markdown",
	})
	if err != nil {
		t.Fatal(err)
	}
	if status != 200 || string(data) != "# Hello" {
		t.Fatalf("status=%d data=%q", status, data)
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
	_, _, err := c.Request(context.Background(), brightdata.RequestOpts{
		Zone: "z",
		URL:  "https://example.com",
	})
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("expected 403 error, got %v", err)
	}
}
