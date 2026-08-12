package tools

import (
	"github.com/datumbridge/bright-data-mcp/internal/mcp"
)

// Register returns all Bright Data MCP tools (official tool list parity).
// Tool visibility: BRIGHTDATA_PRO_MODE=true enables all tools; otherwise Rapid tools only.
// Optional: BRIGHTDATA_GROUPS=ecommerce,social or BRIGHTDATA_TOOLS=tool1,tool2
func Register() ([]mcp.ToolDesc, map[string]mcp.ToolHandler) {
	cfg := LoadServerConfig()
	handlers := map[string]mcp.ToolHandler{}
	var descs []mcp.ToolDesc

	add := func(name, desc string, props map[string]interface{}, required []string, h mcp.ToolHandler) {
		descs = append(descs, mcp.ToolDesc{
			Name:        name,
			Description: desc,
			InputSchema: schema(props, required),
		})
		handlers[name] = h
	}

	registerRapidTools(cfg, add)
	if err := registerWebDataTools(cfg, add); err != nil {
		panic("register web_data tools: " + err.Error())
	}
	registerDatasetSearchTools(cfg, add)
	registerBrowserTools(cfg, add)

	return descs, handlers
}
