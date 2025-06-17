package toolshdl

import (
	"encoding/json"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
	"go.kirha.ai/mcpmux"
)

type parser struct{}

func newParser() *parser {
	return &parser{}
}

func (p *parser) toAPITool(tool tools.Tool) mcpmux.Tool {
	raw, _ := json.Marshal(tool.Parameters)
	var inputSchema mcpmux.InputSchema
	if err := json.Unmarshal(raw, &inputSchema); err != nil {
		inputSchema = mcpmux.InputSchema{}
	}
	return mcpmux.Tool{
		Name:        tool.Identifier,
		Description: tool.Description,
		InputSchema: inputSchema,
	}
}

func (p *parser) toAPITools(tools []tools.Tool) []mcpmux.Tool {
	if tools == nil {
		return []mcpmux.Tool{}
	}
	var apiTools []mcpmux.Tool
	for _, t := range tools {
		apiTools = append(apiTools, p.toAPITool(t))
	}
	return apiTools
}

func (p *parser) toAPIResult(result *tools.ToolExecutionResult) *mcpmux.ToolResult {
	raw, _ := json.Marshal(result.Result)
	var toolResult mcpmux.ToolResult
	if err := json.Unmarshal(raw, &toolResult); err != nil {
		return &mcpmux.ToolResult{
			Content: []mcpmux.Content{
				{
					Type:     "text",
					Text:     "Error parsing tool result",
					MimeType: "text/plain",
				},
			},
			IsError: true,
		}
	}

	return &toolResult
}
