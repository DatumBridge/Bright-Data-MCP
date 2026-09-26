package mcp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	ProtocolVersion = "2024-11-05"
	SessionHeader   = "Mcp-Session-Id"
	MaxBodySize     = 1 << 20 // 1 MiB
	ServerName      = "bright-data-mcp"
	ServerVersion   = "1.0.0"
)

// ToolHandler executes a named MCP tool and returns result content.
type ToolHandler func(args json.RawMessage) (result map[string]interface{})

// ToolDesc is an MCP tools/list descriptor.
type ToolDesc struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Meta        map[string]interface{} `json:"_meta,omitempty"`
}

// CapabilityMeta builds MCP `_meta` consumed by DatumBridge Tool Registry ingest.
func CapabilityMeta(caps ...string) map[string]interface{} {
	values := make([]string, 0, len(caps))
	for _, c := range caps {
		if c != "" {
			values = append(values, c)
		}
	}
	return map[string]interface{}{
		"capabilities": values,
		"datumbridge":  map[string]interface{}{"capabilities": values},
	}
}

// Server is the Streamable HTTP MCP tool-server.
type Server struct {
	Sessions *SessionStore
	Tools    []ToolDesc
	Handlers map[string]ToolHandler
}

type rpcEnvelope struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id"`
}

type callParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// HandleStreamableHTTP implements initialize, tools/list, tools/call.
func (s *Server) HandleStreamableHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
		return
	}
	_ = r.Body.Close()

	var env rpcEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		WriteRPCError(w, nil, -32700, "parse error")
		return
	}
	if env.JSONRPC != "2.0" {
		WriteRPCError(w, env.ID, -32600, "invalid request")
		return
	}

	switch env.Method {
	case "initialize":
		s.handleInitialize(w, env.ID)
	case "notifications/initialized":
		w.WriteHeader(http.StatusAccepted)
	case "tools/list":
		s.handleToolsList(w, r, env.ID)
	case "tools/call":
		s.handleToolsCall(w, r, env)
	default:
		WriteRPCError(w, env.ID, -32601, "method not found: "+env.Method)
	}
}

func (s *Server) handleInitialize(w http.ResponseWriter, id json.RawMessage) {
	sid := s.Sessions.Create()
	w.Header().Set(SessionHeader, sid)
	w.Header().Set("Content-Type", "application/json")
	result := map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]string{
			"name":    ServerName,
			"version": ServerVersion,
		},
	}
	WriteRPCResult(w, id, result)
	preview := sid
	if len(preview) > 8 {
		preview = preview[:8] + "…"
	}
	log.Debug().Str("mcp_session", preview).Msg("mcp initialize")
}

func (s *Server) handleToolsList(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
	if !s.Sessions.Valid(r.Header.Get(SessionHeader)) {
		WriteRPCError(w, id, -32000, "invalid or missing Mcp-Session-Id (call initialize first)")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	WriteRPCResult(w, id, map[string]interface{}{"tools": s.Tools})
}

func (s *Server) handleToolsCall(w http.ResponseWriter, r *http.Request, env rpcEnvelope) {
	if !s.Sessions.Valid(r.Header.Get(SessionHeader)) {
		WriteRPCError(w, env.ID, -32000, "invalid or missing Mcp-Session-Id (call initialize first)")
		return
	}

	var params callParams
	if len(env.Params) > 0 && string(env.Params) != "null" {
		if err := json.Unmarshal(env.Params, &params); err != nil {
			WriteRPCError(w, env.ID, -32602, "invalid params")
			return
		}
	}
	if params.Arguments == nil {
		params.Arguments = []byte("{}")
	}

	w.Header().Set("Content-Type", "application/json")
	h, ok := s.Handlers[params.Name]
	if !ok {
		WriteRPCError(w, env.ID, -32601, "unknown tool: "+params.Name)
		return
	}
	WriteRPCResult(w, env.ID, h(params.Arguments))
}

// ToolResultText builds a successful MCP tool result.
func ToolResultText(text string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]string{{"type": "text", "text": text}},
		"isError": false,
	}
}

// ToolResultError builds an MCP tool domain error (isError=true).
func ToolResultError(text string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]string{{"type": "text", "text": text}},
		"isError": true,
	}
}

// WriteRPCResult writes a JSON-RPC success response.
func WriteRPCResult(w http.ResponseWriter, id json.RawMessage, result interface{}) {
	out := map[string]interface{}{"jsonrpc": "2.0", "result": result}
	if len(id) > 0 && string(id) != "null" {
		var idVal interface{}
		if err := json.Unmarshal(id, &idVal); err == nil {
			out["id"] = idVal
		}
	}
	_ = json.NewEncoder(w).Encode(out)
}

// WriteRPCError writes a JSON-RPC error response.
func WriteRPCError(w http.ResponseWriter, id json.RawMessage, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	out := map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	if len(id) > 0 && string(id) != "null" {
		var idVal interface{}
		if err := json.Unmarshal(id, &idVal); err == nil {
			out["id"] = idVal
		}
	}
	_ = json.NewEncoder(w).Encode(out)
}
