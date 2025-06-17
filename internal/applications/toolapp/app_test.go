package toolapp

import (
	"context"
	"errors"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
	"testing"
	"time"

	domainErrors "go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"go.kirha.ai/mcpmux"
)

// mockKirhaClient is a mock implementation of the KirhaClient interface for testing.
type mockKirhaClient struct {
	tools      []tools.Tool
	toolsError error
	execResult *tools.ToolExecutionResult
	execError  error
}

// ListTools returns the mocked tools list or error.
func (m *mockKirhaClient) ListTools(ctx context.Context) ([]tools.Tool, error) {
	return m.tools, m.toolsError
}

// ExecuteTool returns the mocked execution result or error.
func (m *mockKirhaClient) ExecuteTool(ctx context.Context, name string, arguments mcpmux.Args) (*tools.ToolExecutionResult, error) {
	return m.execResult, m.execError
}

func TestMCPGatewayApplication_ListTools(t *testing.T) {
	tests := []struct {
		name          string
		mockTools     []tools.Tool
		mockError     error
		expectedCount int
		expectError   bool
	}{
		{
			name: "successful list tools - single tool",
			mockTools: []tools.Tool{
				{
					ID:          "1",
					Identifier:  "test-tool-1",
					Name:        "Test Tool 1",
					Description: "A test tool",
					MCPID:       "mcp-1",
					VerticalIDs: []string{"vertical-1"},
					Parameters:  map[string]interface{}{"type": "object"},
					Outputs:     map[string]interface{}{"type": "string"},
				},
			},
			mockError:     nil,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "successful list tools - multiple tools",
			mockTools: []tools.Tool{
				{
					ID:          "1",
					Identifier:  "test-tool-1",
					Name:        "Test Tool 1",
					Description: "First test tool",
					MCPID:       "mcp-1",
					VerticalIDs: []string{"vertical-1"},
				},
				{
					ID:          "2",
					Identifier:  "test-tool-2",
					Name:        "Test Tool 2",
					Description: "Second test tool",
					MCPID:       "mcp-2",
					VerticalIDs: []string{"vertical-1", "vertical-2"},
				},
			},
			mockError:     nil,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "successful list tools - empty list",
			mockTools:     []tools.Tool{},
			mockError:     nil,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "client error - network timeout",
			mockTools:     nil,
			mockError:     domainErrors.ErrNetworkTimeout,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "client error - unauthorized",
			mockTools:     nil,
			mockError:     domainErrors.ErrUnauthorized,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "client error - generic error",
			mockTools:     nil,
			mockError:     errors.New("unexpected error"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockKirhaClient{
				tools:      tt.mockTools,
				toolsError: tt.mockError,
			}

			app := New(mockClient)
			tools, err := app.ListTools(context.Background())

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(tools) != tt.expectedCount {
				t.Errorf("expected %d tools, got %d", tt.expectedCount, len(tools))
			}
		})
	}
}

func TestMCPGatewayApplication_ExecuteTool(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		arguments   mcpmux.Args
		mockResult  *tools.ToolExecutionResult
		mockError   error
		expectError bool
		errorType   error
	}{
		{
			name:      "successful execution - simple result",
			toolName:  "test-tool",
			arguments: mcpmux.Args{"param": "value"},
			mockResult: &tools.ToolExecutionResult{
				ToolName:  "test-tool",
				Success:   true,
				Result:    map[string]interface{}{"output": "success"},
				Duration:  time.Millisecond * 100,
				Timestamp: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:      "successful execution - complex result",
			toolName:  "data-processor",
			arguments: mcpmux.Args{"input": []string{"a", "b", "c"}, "format": "json"},
			mockResult: &tools.ToolExecutionResult{
				ToolName: "data-processor",
				Success:  true,
				Result: map[string]interface{}{
					"processed": []interface{}{"A", "B", "C"},
					"count":     3,
					"metadata":  map[string]interface{}{"format": "json", "version": "1.0"},
				},
				Duration:  time.Millisecond * 250,
				Timestamp: time.Now(),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:      "execution failure - tool reports failure",
			toolName:  "test-tool",
			arguments: mcpmux.Args{"invalid": "param"},
			mockResult: &tools.ToolExecutionResult{
				ToolName:  "test-tool",
				Success:   false,
				Error:     "Invalid parameters provided",
				Duration:  time.Millisecond * 50,
				Timestamp: time.Now(),
			},
			mockError:   nil,
			expectError: true,
			errorType:   domainErrors.ErrToolExecutionFailed,
		},
		{
			name:        "client error - tool not found",
			toolName:    "nonexistent-tool",
			arguments:   mcpmux.Args{},
			mockResult:  nil,
			mockError:   domainErrors.ErrToolNotFound,
			expectError: true,
			errorType:   domainErrors.ErrToolNotFound,
		},
		{
			name:        "client error - unauthorized",
			toolName:    "test-tool",
			arguments:   mcpmux.Args{},
			mockResult:  nil,
			mockError:   domainErrors.ErrUnauthorized,
			expectError: true,
			errorType:   domainErrors.ErrUnauthorized,
		},
		{
			name:        "client error - network timeout",
			toolName:    "test-tool",
			arguments:   mcpmux.Args{"timeout": true},
			mockResult:  nil,
			mockError:   domainErrors.ErrNetworkTimeout,
			expectError: true,
			errorType:   domainErrors.ErrNetworkTimeout,
		},
		{
			name:        "client error - generic error",
			toolName:    "test-tool",
			arguments:   mcpmux.Args{},
			mockResult:  nil,
			mockError:   errors.New("unexpected error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockKirhaClient{
				execResult: tt.mockResult,
				execError:  tt.mockError,
			}

			app := New(mockClient)
			result, err := app.ExecuteTool(context.Background(), tt.toolName, tt.arguments)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check for specific error type if specified
			if tt.expectError && tt.errorType != nil && !errors.Is(err, tt.errorType) {
				t.Errorf("expected error type %v, got %v", tt.errorType, err)
			}

			if !tt.expectError && result == nil {
				t.Error("expected result but got nil")
			}

			// Verify result properties for successful executions
			if result != nil && !tt.expectError {
				if result.ToolName != tt.toolName {
					t.Errorf("expected tool name %s, got %s", tt.toolName, result.ToolName)
				}
				if !result.Success {
					t.Error("expected successful execution but got failure")
				}
			}
		})
	}
}
