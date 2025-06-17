package toolshdl

import (
	"context"
	"go.kirha.ai/kirha-mcp-gateway/internal/applications/toolapp"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"log/slog"

	"go.kirha.ai/mcpmux"
)

// Handler implements the mcpmux.ToolAggregator interface for MCP tool operations.
// It translates between the MCP protocol and the application layer.
type Handler struct {
	app    *toolapp.Application
	parser *parser
	logger *slog.Logger
}

// New creates a new MCP tools handler with the provided tool application.
// The handler will use the application to process tool operations and a parser
// to convert between domain and MCP protocol types.
//
// Parameters:
//   - app: The application layer for tool management operations
//
// Returns:
//   - *Handler: A new handler instance implementing mcpmux.ToolAggregator
func New(app *toolapp.Application) *Handler {
	return &Handler{
		app:    app,
		parser: newParser(),
		logger: logger.New("mcp_handler"),
	}
}

// ListTools handles MCP tool listing requests.
// It retrieves available tools from the application layer and converts them to MCP format.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//
// Returns:
//   - []mcpmux.Tool: List of tools in MCP protocol format
//   - error: An MCP protocol error if the operation fails
func (h *Handler) ListTools(ctx context.Context) ([]mcpmux.Tool, error) {
	h.logger.InfoContext(ctx, "handling list tools request")

	tools, err := h.app.ListTools(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list tools", slog.String("error", err.Error()))
		return nil, mcpmux.NewError(mcpmux.ErrorCodeInternal, "failed to list tools", err)
	}

	h.logger.InfoContext(ctx, "successfully handled list tools request", slog.Int("count", len(tools)))
	return h.parser.toAPITools(tools), nil
}

// ExecuteTool handles MCP tool execution requests.
// It executes the specified tool through the application layer and converts the result to MCP format.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//   - name: The name of the tool to execute
//   - arguments: The arguments to pass to the tool in MCP format
//
// Returns:
//   - *mcpmux.ToolResult: The tool execution result in MCP protocol format
//   - error: An MCP protocol error if the operation fails
func (h *Handler) ExecuteTool(ctx context.Context, name string, arguments mcpmux.Args) (*mcpmux.ToolResult, error) {
	h.logger.InfoContext(ctx, "handling execute tool request", slog.String("tool_name", name))

	result, err := h.app.ExecuteTool(ctx, name, arguments)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to execute tool", slog.String("tool_name", name), slog.String("error", err.Error()))
		return nil, mcpmux.NewError(mcpmux.ErrorCodeInternal, "failed to execute tool", err)
	}

	h.logger.InfoContext(ctx, "successfully handled execute tool request", slog.String("tool_name", name))
	return h.parser.toAPIResult(result), nil
}
