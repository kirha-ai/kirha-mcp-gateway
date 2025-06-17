package tools

import "time"

// Tool represents a tool that can be executed within the MCP.
type Tool struct {
	ID          string
	Identifier  string
	Name        string
	Description string
	MCPID       string
	VerticalIDs []string
	Parameters  interface{}
	Outputs     interface{}
}

// ToolExecutionResult represents the result of executing a tool.
type ToolExecutionResult struct {
	ToolName  string
	Result    map[string]interface{}
	Duration  time.Duration
	Success   bool
	Error     string
	Timestamp time.Time
}
