package clients

import (
	"context"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
	"go.kirha.ai/mcpmux"
)

// KirhaClient defines the interface for interacting with the Kirha API to manage tools.
type KirhaClient interface {
	// ListTools retrieves the list of tools available in the Kirha API.
	ListTools(ctx context.Context) ([]tools.Tool, error)
	// ExecuteTool executes a tool with the given name and arguments.
	ExecuteTool(ctx context.Context, name string, arguments mcpmux.Args) (*tools.ToolExecutionResult, error)
}
