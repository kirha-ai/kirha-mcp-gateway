package toolapp

import (
	"context"
	"fmt"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/ports/clients"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"log/slog"

	"go.kirha.ai/mcpmux"
)

// Application implements the tool management use cases.
// It orchestrates the business logic for listing and executing tools through the Kirha API.
type Application struct {
	kirhaClient clients.KirhaClient
	logger      *slog.Logger
}

// New creates a new Application instance with the provided Kirha client.
// The application will use this client to interact with the Kirha API for all tool operations.
//
// Parameters:
//   - kirhaClient: The client implementation for communicating with the Kirha API
//
// Returns:
//   - *Application: A new application instance ready to handle tool operations
func New(kirhaClient clients.KirhaClient) *Application {
	return &Application{
		kirhaClient: kirhaClient,
		logger:      logger.New("gateway_application"),
	}
}

// ListTools retrieves all available tools from the Kirha API.
// It delegates to the Kirha client and handles logging of the operation.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//
// Returns:
//   - []tools.Tool: List of available tools
//   - error: An error if the operation fails
func (a *Application) ListTools(ctx context.Context) ([]tools.Tool, error) {
	a.logger.InfoContext(ctx, "listing tools from Kirha API")

	domainTools, err := a.kirhaClient.ListTools(ctx)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to list tools from Kirha API", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	a.logger.InfoContext(ctx, "successfully listed tools from Kirha API", slog.Int("count", len(domainTools)))
	return domainTools, nil
}

// ExecuteTool executes a specific tool by name with the provided arguments.
// It delegates to the Kirha client for execution and validates the result.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//   - name: The name of the tool to execute
//   - arguments: The arguments to pass to the tool as mcpmux.Args
//
// Returns:
//   - *tools.ToolExecutionResult: The result of the tool execution
//   - error: An error if the tool execution fails or if the tool reports failure
//
// Errors:
//   - Returns the underlying error from the Kirha client if the execution fails
//   - Returns errors.ErrToolExecutionFailed if the tool execution reports failure in the result
func (a *Application) ExecuteTool(ctx context.Context, name string, arguments mcpmux.Args) (*tools.ToolExecutionResult, error) {
	a.logger.InfoContext(ctx, "executing tool", slog.String("tool_name", name))

	result, err := a.kirhaClient.ExecuteTool(ctx, name, arguments)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to execute tool", slog.String("tool_name", name), slog.String("error", err.Error()))
		return nil, err
	}

	if !result.Success {
		a.logger.ErrorContext(ctx, "tool execution failed", slog.String("tool_name", name), slog.String("error", result.Error))
		return nil, errors.ErrToolExecutionFailed
	}

	a.logger.InfoContext(ctx, "tool executed successfully", slog.String("tool_name", name), slog.Duration("duration", result.Duration))
	return result, nil
}
