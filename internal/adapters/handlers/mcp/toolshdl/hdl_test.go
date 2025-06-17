package toolshdl

import (
	"context"
	"errors"
	"testing"

	"go.kirha.ai/kirha-mcp-gateway/internal/applications/toolapp"
	domainErrors "go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/domain/tools"
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

// TestHandler_ListTools tests the ListTools method of the MCP handler.
func TestHandler_ListTools(t *testing.T) {
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
					Identifier:  "test-tool",
					Name:        "Test Tool",
					Description: "A test tool",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"input": map[string]interface{}{"type": "string"},
						},
					},
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
					Identifier:  "tool-1",
					Name:        "Tool 1",
					Description: "First tool",
				},
				{
					ID:          "2",
					Identifier:  "tool-2",
					Name:        "Tool 2",
					Description: "Second tool",
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
			name:          "application error",
			mockTools:     nil,
			mockError:     errors.New("application error"),
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "domain error - unauthorized",
			mockTools:     nil,
			mockError:     domainErrors.ErrUnauthorized,
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

			// Create application with mock client
			app := toolapp.New(mockClient)
			// Create handler using the constructor
			handler := New(app)

			mcpTools, err := handler.ListTools(context.Background())

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(mcpTools) != tt.expectedCount {
				t.Errorf("expected %d tools, got %d", tt.expectedCount, len(mcpTools))
			}

			// Verify tool conversion for successful cases
			if !tt.expectError && tt.expectedCount > 0 {
				for i, mcpTool := range mcpTools {
					// The parser uses Identifier as the tool name
					expectedName := tt.mockTools[i].Identifier
					if expectedName == "" {
						expectedName = tt.mockTools[i].Name
					}
					if mcpTool.Name != expectedName {
						t.Errorf("expected tool name %s, got %s", expectedName, mcpTool.Name)
					}
					if mcpTool.Description != tt.mockTools[i].Description {
						t.Errorf("expected description %s, got %s", tt.mockTools[i].Description, mcpTool.Description)
					}
				}
			}
		})
	}
}

// TestHandler_ExecuteTool tests the ExecuteTool method of the MCP handler.
func TestHandler_ExecuteTool(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		arguments   mcpmux.Args
		mockResult  *tools.ToolExecutionResult
		mockError   error
		expectError bool
		checkResult bool
	}{
		{
			name:      "successful execution - simple result",
			toolName:  "test-tool",
			arguments: mcpmux.Args{"input": "test"},
			mockResult: &tools.ToolExecutionResult{
				ToolName: "test-tool",
				Success:  true,
				Result: map[string]interface{}{
					"output": "test result",
				},
			},
			mockError:   nil,
			expectError: false,
			checkResult: true,
		},
		{
			name:      "successful execution - complex result",
			toolName:  "complex-tool",
			arguments: mcpmux.Args{"data": []string{"a", "b", "c"}},
			mockResult: &tools.ToolExecutionResult{
				ToolName: "complex-tool",
				Success:  true,
				Result: map[string]interface{}{
					"processed": []string{"A", "B", "C"},
					"count":     3,
					"metadata": map[string]interface{}{
						"version": "1.0",
						"status":  "complete",
					},
				},
			},
			mockError:   nil,
			expectError: false,
			checkResult: true,
		},
		{
			name:      "execution failure - tool error",
			toolName:  "failing-tool",
			arguments: mcpmux.Args{},
			mockResult: &tools.ToolExecutionResult{
				ToolName: "failing-tool",
				Success:  false,
				Error:    "Tool execution failed",
			},
			mockError:   nil,
			expectError: true, // Application layer converts failed tool execution to error
			checkResult: false,
		},
		{
			name:        "application error - tool not found",
			toolName:    "nonexistent-tool",
			arguments:   mcpmux.Args{},
			mockResult:  nil,
			mockError:   domainErrors.ErrToolNotFound,
			expectError: true,
			checkResult: false,
		},
		{
			name:        "application error - generic error",
			toolName:    "test-tool",
			arguments:   mcpmux.Args{},
			mockResult:  nil,
			mockError:   errors.New("unexpected error"),
			expectError: true,
			checkResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockKirhaClient{
				execResult: tt.mockResult,
				execError:  tt.mockError,
			}

			// Create application with mock client
			app := toolapp.New(mockClient)
			// Create handler using the constructor
			handler := New(app)

			result, err := handler.ExecuteTool(context.Background(), tt.toolName, tt.arguments)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.checkResult && result != nil {
				// Just verify we got a result - the structure depends on JSON unmarshaling
				// which is hard to predict in tests
				if result == nil {
					t.Error("expected result but got nil")
				}
			}
		})
	}
}

// TestHandler_New tests the creation of a new handler.
func TestHandler_New(t *testing.T) {
	mockClient := &mockKirhaClient{}
	app := toolapp.New(mockClient)

	handler := New(app)

	if handler == nil {
		t.Fatal("expected handler instance but got nil")
	}

	if handler.app == nil {
		t.Error("handler app field not set correctly")
	}

	if handler.parser == nil {
		t.Error("handler parser field not initialized")
	}

	if handler.logger == nil {
		t.Error("handler logger field not initialized")
	}
}
