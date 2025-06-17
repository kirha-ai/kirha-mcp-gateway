package kirha

type ListToolsResponse struct {
	Tools      []ToolResponse `json:"tools"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

type ToolResponse struct {
	ID          string      `json:"id"`
	Identifier  string      `json:"identifier"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	MCPID       string      `json:"mcp_id"`
	VerticalIDs []string    `json:"vertical_ids,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
	Outputs     interface{} `json:"outputs,omitempty"`
}

type ExecuteToolRequest struct {
	Arguments map[string]interface{} `json:"arguments"`
}

type ExecuteToolResponse struct {
	Result map[string]interface{} `json:"result"`
}
