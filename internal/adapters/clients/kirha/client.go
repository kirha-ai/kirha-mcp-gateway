package kirha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/ports/clients"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"log/slog"
	"net/http"
	"time"

	"go.kirha.ai/mcpmux"
)

// client implements the KirhaClient interface using HTTP communication.
type client struct {
	config     *Config
	httpClient *http.Client
	logger     *slog.Logger
}

// New creates a new instance of the HTTP Kirha client with the provided configuration.
// The client is configured with the timeout specified in the configuration.
//
// Parameters:
//   - config: Configuration containing base URL, API key, timeout, and vertical ID
//
// Returns:
//   - clients.KirhaClient: An implementation of the KirhaClient interface
func New(config *Config) clients.KirhaClient {
	return &client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger.New("kirha_client"),
	}
}

// ListTools retrieves all available tools from the Kirha API for the configured vertical.
// It makes an HTTP GET request to the /mcp/v1/tools endpoint with the vertical_id filter.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//
// Returns:
//   - []tools.Tool: List of available tools for the configured vertical
//   - error: An error if the operation fails
//
// Errors:
//   - errors.ErrUnauthorized: If the API key is invalid or missing
//   - errors.ErrNetworkTimeout: If the request times out
//   - errors.ErrInvalidResponse: If the response cannot be decoded
//   - errors.ErrInternalServer: For other HTTP errors (status >= 400)
func (c *client) ListTools(ctx context.Context) ([]tools.Tool, error) {
	url := fmt.Sprintf("%s/mcp/v1/tools?limit=-1&vertical_id=%s", c.config.BaseURL, c.config.VerticalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to create request", slog.String("error", err.Error()))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to execute request", slog.String("error", err.Error()))
		return nil, errors.ErrNetworkTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.logger.ErrorContext(ctx, "unauthorized request", "status_code", resp.StatusCode)
		return nil, errors.ErrUnauthorized
	}

	if resp.StatusCode >= 400 {
		c.logger.ErrorContext(ctx, "API error", slog.Int("status_code", resp.StatusCode))
		return nil, errors.ErrInternalServer
	}

	var response ListToolsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.ErrorContext(ctx, "failed to decode response", slog.String("error", err.Error()))
		return nil, errors.ErrInvalidResponse
	}

	listedTools := make([]tools.Tool, len(response.Tools))
	for i, tool := range response.Tools {
		listedTools[i] = tools.Tool{
			ID:          tool.ID,
			Identifier:  tool.Identifier,
			Name:        tool.Name,
			Description: tool.Description,
			MCPID:       tool.MCPID,
			VerticalIDs: tool.VerticalIDs,
			Parameters:  tool.Parameters,
			Outputs:     tool.Outputs,
		}
	}

	c.logger.InfoContext(ctx, "successfully listed tools", slog.Int("count", len(listedTools)))
	return listedTools, nil
}

// ExecuteTool executes a specific tool by name with the provided arguments.
// It makes an HTTP POST request to the /mcp/v1/tools/{name}/execute endpoint.
// The method always returns a ToolExecutionResult, even on failure, to provide execution metadata.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and request tracing
//   - name: The name of the tool to execute
//   - arguments: The arguments to pass to the tool as mcpmux.Args
//
// Returns:
//   - *tools.ToolExecutionResult: The result of the tool execution, including success status and metadata
//   - error: An error if the operation fails
//
// Errors:
//   - errors.ErrUnauthorized: If the API key is invalid or missing
//   - errors.ErrToolNotFound: If the specified tool does not exist
//   - errors.ErrNetworkTimeout: If the request times out
//   - errors.ErrInvalidResponse: If the response cannot be decoded
//   - errors.ErrToolExecutionFailed: For other HTTP errors (status >= 400)
//
// Note: Even when an error is returned, the ToolExecutionResult will contain
// partial information about the failed execution, including duration and error details.
func (c *client) ExecuteTool(ctx context.Context, name string, arguments mcpmux.Args) (*tools.ToolExecutionResult, error) {
	startTime := time.Now()

	url := fmt.Sprintf("%s/mcp/v1/tools/%s/execute", c.config.BaseURL, name)

	reqBody := ExecuteToolRequest{
		Arguments: arguments,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to marshal request body", slog.String("error", err.Error()), slog.String("tool", name))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to create request", slog.String("error", err.Error()), slog.String("tool", name))
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to execute request", slog.String("error", err.Error()), "tool", name)
		return &tools.ToolExecutionResult{
			ToolName:  name,
			Success:   false,
			Error:     err.Error(),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, errors.ErrNetworkTimeout
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logger.Error("unauthorized request", "status_code", resp.StatusCode, "tool", name)
		return &tools.ToolExecutionResult{
			ToolName:  name,
			Success:   false,
			Error:     "unauthorized",
			Duration:  duration,
			Timestamp: time.Now(),
		}, errors.ErrUnauthorized
	}

	if resp.StatusCode == http.StatusNotFound {
		c.logger.Error("tool not found", "tool", name)
		return &tools.ToolExecutionResult{
			ToolName:  name,
			Success:   false,
			Error:     "tool not found",
			Duration:  duration,
			Timestamp: time.Now(),
		}, errors.ErrToolNotFound
	}

	if resp.StatusCode >= 400 {
		c.logger.Error("API error", "status_code", resp.StatusCode, "tool", name)
		return &tools.ToolExecutionResult{
			ToolName:  name,
			Success:   false,
			Error:     fmt.Sprintf("API error: %d", resp.StatusCode),
			Duration:  duration,
			Timestamp: time.Now(),
		}, errors.ErrToolExecutionFailed
	}

	var response ExecuteToolResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("failed to decode response", slog.String("error", err.Error()), "tool", name)
		return &tools.ToolExecutionResult{
			ToolName:  name,
			Success:   false,
			Error:     "invalid response",
			Duration:  duration,
			Timestamp: time.Now(),
		}, errors.ErrInvalidResponse
	}

	c.logger.Info("successfully executed tool", "tool", name, "duration", duration)
	return &tools.ToolExecutionResult{
		ToolName:  name,
		Result:    response.Result,
		Success:   true,
		Duration:  duration,
		Timestamp: time.Now(),
	}, nil
}
