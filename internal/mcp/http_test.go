package mcp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/mcp"
	"github.com/datumbridge/bright-data-mcp/internal/tools"
)

func newTestServer() *mcp.Server {
	descs, handlers := tools.Register()
	return &mcp.Server{
		Sessions: mcp.NewSessionStore(),
		Tools:    descs,
		Handlers: handlers,
	}
}

func postRPC(t *testing.T, s *mcp.Server, body map[string]interface{}, session string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if session != "" {
		req.Header.Set(mcp.SessionHeader, session)
	}
	rr := httptest.NewRecorder()
	s.HandleStreamableHTTP(rr, req)
	return rr
}

func TestInitializeAndToolsList(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params":  map[string]interface{}{},
	}, "")
	if rr.Code != 200 {
		t.Fatalf("status %d", rr.Code)
	}
	sid := rr.Header().Get(mcp.SessionHeader)
	if sid == "" {
		t.Fatal("expected Mcp-Session-Id")
	}

	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}, sid)
	var resp map[string]interface{}
	if err := json.Unmarshal(rr2.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("no result: %s", rr2.Body.String())
	}
	toolsList, ok := result["tools"].([]interface{})
	if !ok || len(toolsList) < 7 {
		t.Fatalf("expected at least 7 rapid tools, got %v", result["tools"])
	}
	capsByName := map[string][]string{}
	for _, raw := range toolsList {
		tool, _ := raw.(map[string]interface{})
		name, _ := tool["name"].(string)
		meta, _ := tool["_meta"].(map[string]interface{})
		if meta == nil {
			t.Fatalf("tool %s missing _meta", name)
		}
		rawCaps, _ := meta["capabilities"].([]interface{})
		if len(rawCaps) == 0 {
			t.Fatalf("tool %s missing _meta.capabilities", name)
		}
		caps := make([]string, 0, len(rawCaps))
		for _, c := range rawCaps {
			s, _ := c.(string)
			if s != "" {
				caps = append(caps, s)
			}
		}
		capsByName[name] = caps
	}
	if join(capsByName["search_engine"]) == join(capsByName["scrape_as_markdown"]) {
		t.Fatalf("expected distinct capabilities, got %v", capsByName["search_engine"])
	}
}

func join(in []string) string {
	out := ""
	for i, s := range in {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

func TestToolsListRejectsMissingSession(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	}, "")
	var resp map[string]interface{}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] == nil {
		t.Fatalf("expected error: %s", rr.Body.String())
	}
}

func TestUnknownTool(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
	}, "")
	sid := rr.Header().Get(mcp.SessionHeader)
	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      "nope",
			"arguments": map[string]interface{}{},
		},
	}, sid)
	var resp map[string]interface{}
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp)
	errObj, _ := resp["error"].(map[string]interface{})
	if errObj == nil {
		t.Fatalf("expected rpc error: %s", rr2.Body.String())
	}
	code, _ := errObj["code"].(float64)
	if int(code) != -32601 {
		t.Fatalf("code=%v", code)
	}
}
