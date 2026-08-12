package tools

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
	"github.com/datumbridge/bright-data-mcp/internal/mcp"
)

var (
	catalogOnce sync.Once
	catalogData []brightdata.DatasetTool
	catalogErr  error
)

func loadDatasetCatalogCached() ([]brightdata.DatasetTool, error) {
	catalogOnce.Do(func() {
		catalogData, catalogErr = brightdata.LoadDatasetCatalog()
	})
	return catalogData, catalogErr
}

func registerWebDataTools(cfg ServerConfig, add toolAdder) error {
	catalog, err := loadDatasetCatalogCached()
	if err != nil {
		return err
	}
	for _, tool := range catalog {
		name := "web_data_" + tool.ID
		if !cfg.toolEnabled(name) {
			continue
		}
		props := baseProps(nil)
		required := []string{}
		for _, input := range tool.Inputs {
			desc := input
			if input == "url" {
				desc = "Target URL (http or https)"
			}
			if input == "prompt" {
				desc = "Prompt for AI insights"
			}
			if input == "package_name" {
				desc = "Package name (e.g. @scope/pkg or langchain-brightdata)"
			}
			props[input] = map[string]interface{}{
				"type":        "string",
				"description": desc,
			}
			if def, ok := tool.Defaults[input]; ok {
				props[input].(map[string]interface{})["default"] = def
			}
			if _, hasDef := tool.Defaults[input]; !hasDef && input != "start_date" && input != "end_date" && input != "days_back" {
				required = append(required, input)
			}
		}
		toolCopy := tool
		add(name, tool.Description, props, required, func(raw json.RawMessage) map[string]interface{} {
			recordToolCall(name)
			return handleWebDataTool(toolCopy, raw)
		})
	}
	return nil
}

func handleWebDataTool(tool brightdata.DatasetTool, raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	input := map[string]string{}
	for _, key := range tool.Inputs {
		if v := strArg(m, key); v != "" {
			input[key] = v
		}
	}
	cctx, cancel := pollCtx()
	defer cancel()
	data, err := client.TriggerWebData(cctx, tool, input)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(string(data))
}

// searchableDatasets from official MCP.
var searchableDatasets = []string{
	"gd_l1viktl72bvl7bjuj0",
	"gd_me5ppxjr2ge6icjuh0",
	"gd_l1vikfnt1wgvvqz95w",
}

func registerDatasetSearchTools(cfg ServerConfig, add toolAdder) {
	if cfg.toolEnabled("list_dataset_fields") {
		add("list_dataset_fields",
			"List filterable fields of a searchable Bright Data dataset. Supported dataset_id: "+
				strings.Join(searchableDatasets, ", "),
			baseProps(map[string]interface{}{
				"dataset_id": map[string]interface{}{
					"type":        "string",
					"description": "Dataset ID",
					"enum":        searchableDatasets,
				},
			}),
			[]string{"dataset_id"},
			func(raw json.RawMessage) map[string]interface{} {
				recordToolCall("list_dataset_fields")
				return handleListDatasetFields(raw)
			},
		)
	}
	if cfg.toolEnabled("search_dataset") {
		add("search_dataset",
			"Search a Bright Data dataset by filter (Elasticsearch-backed). Use list_dataset_fields first.",
			baseProps(map[string]interface{}{
				"dataset_id": map[string]interface{}{
					"type": "string", "enum": searchableDatasets,
				},
				"filter": map[string]interface{}{
					"type":        "object",
					"description": "Filter tree: group {operator, filters} or leaf {name, operator, value}",
				},
				"size": map[string]interface{}{"type": "integer"},
				"sort": map[string]interface{}{"type": "array"},
			}),
			[]string{"dataset_id", "filter"},
			func(raw json.RawMessage) map[string]interface{} {
				recordToolCall("search_dataset")
				return handleSearchDataset(raw)
			},
		)
	}
}

func handleListDatasetFields(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	datasetID := strArg(m, "dataset_id")
	if datasetID == "" {
		return mcp.ToolResultError("dataset_id is required")
	}
	cctx, cancel := pollCtx()
	defer cancel()
	data, err := client.DatasetMetadata(cctx, datasetID)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(string(data))
}

func handleSearchDataset(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	datasetID := strArg(m, "dataset_id")
	if datasetID == "" {
		return mcp.ToolResultError("dataset_id is required")
	}
	body := map[string]interface{}{}
	if f, ok := m["filter"]; ok {
		body["filter"] = f
	}
	if sz := intArg(m, "size", 0); sz > 0 {
		body["size"] = sz
	}
	if s, ok := m["sort"]; ok {
		body["sort"] = s
	}
	cctx, cancel := pollCtx()
	defer cancel()
	data, err := client.SearchDataset(cctx, datasetID, body)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(string(data))
}

func untrustedTextResult(body string) map[string]interface{} {
	return mcp.ToolResultText(untrustedPrefix + body)
}
