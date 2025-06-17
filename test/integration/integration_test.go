// Package integration provides integration tests for the Kirha MCP Gateway.
// These tests verify the full flow from HTTP client through application layer.
package integration

import (
	"context"
	"encoding/json"
	"go.kirha.ai/kirha-mcp-gateway/internal/adapters/clients/kirha"
	"go.kirha.ai/kirha-mcp-gateway/internal/applications/toolapp"
	"go.kirha.ai/mcpmux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestFullIntegration tests the complete flow of listing and executing tools.
func TestFullIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/mcp/v1/tools":
			response := kirha.ListToolsResponse{
				Tools: []kirha.ToolResponse{
					{
						ID:          "1",
						Identifier:  "integration-test-tool",
						Name:        "Integration Test Tool",
						Description: "A tool for integration testing",
						MCPID:       "mcp-1",
						VerticalIDs: []string{"test-vertical"},
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		case "/mcp/v1/tools/integration-test-tool/execute":
			response := kirha.ExecuteToolResponse{
				Result: map[string]interface{}{
					"text": "Integration test successful",
					"data": "test data",
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &kirha.Config{
		APIKey:     "test-api-key",
		VerticalID: "test-vertical",
		BaseURL:    server.URL,
		Timeout:    time.Second * 5,
	}

	client := kirha.New(config)
	app := toolapp.New(client)

	ctx := context.Background()

	// Test listing tools
	tools, err := app.ListTools(ctx)
	if err != nil {
		t.Fatalf("failed to list tools: %v", err)
	}

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}

	tool := tools[0]
	if tool.Name != "Integration Test Tool" {
		t.Errorf("expected tool name 'Integration Test Tool', got '%s'", tool.Name)
	}

	// Test executing tool
	result, err := app.ExecuteTool(ctx, "integration-test-tool", mcpmux.Args{
		"param": "value",
	})
	if err != nil {
		t.Fatalf("failed to execute tool: %v", err)
	}

	if result == nil {
		t.Fatal("expected result but got nil")
	}

	if !result.Success {
		t.Fatal("expected successful execution")
	}

	if result.Result == nil {
		t.Fatal("expected result data but got nil")
	}

	// Verify result contains expected data
	if text, ok := result.Result["text"].(string); !ok || text != "Integration test successful" {
		t.Errorf("unexpected result text: %v", result.Result["text"])
	}

	t.Logf("Integration test completed successfully")
}

// TestErrorHandling tests error scenarios in the integration flow.
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		toolName      string
		expectListErr bool
		expectExecErr bool
	}{
		{
			name: "unauthorized access",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
				}))
			},
			expectListErr: true,
			expectExecErr: true,
		},
		{
			name: "tool not found",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/mcp/v1/tools":
						response := kirha.ListToolsResponse{
							Tools: []kirha.ToolResponse{},
						}
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(response)
					case "/mcp/v1/tools/nonexistent/execute":
						w.WriteHeader(http.StatusNotFound)
						json.NewEncoder(w).Encode(map[string]string{"error": "tool not found"})
					default:
						http.NotFound(w, r)
					}
				}))
			},
			toolName:      "nonexistent",
			expectListErr: false,
			expectExecErr: true,
		},
		{
			name: "server error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
				}))
			},
			expectListErr: true,
			expectExecErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			config := &kirha.Config{
				APIKey:     "test-api-key",
				VerticalID: "test-vertical",
				BaseURL:    server.URL,
				Timeout:    time.Second * 5,
			}

			client := kirha.New(config)
			app := toolapp.New(client)

			ctx := context.Background()

			// Test listing tools
			tools, err := app.ListTools(ctx)
			if tt.expectListErr && err == nil {
				t.Error("expected error when listing tools but got none")
			}
			if !tt.expectListErr && err != nil {
				t.Errorf("unexpected error listing tools: %v", err)
			}

			// Test executing tool if listing was successful or if we're testing specific execution errors
			if !tt.expectListErr || tt.toolName != "" {
				toolName := tt.toolName
				if toolName == "" && len(tools) > 0 {
					toolName = tools[0].Name
				}

				_, err := app.ExecuteTool(ctx, toolName, mcpmux.Args{"test": "param"})
				if tt.expectExecErr && err == nil {
					t.Error("expected error when executing tool but got none")
				}
				if !tt.expectExecErr && err != nil {
					t.Errorf("unexpected error executing tool: %v", err)
				}
			}
		})
	}
}

// TestConcurrentOperations tests concurrent tool listing and execution.
func TestConcurrentOperations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)

		switch r.URL.Path {
		case "/mcp/v1/tools":
			response := kirha.ListToolsResponse{
				Tools: []kirha.ToolResponse{
					{
						ID:          "1",
						Identifier:  "concurrent-tool",
						Name:        "Concurrent Tool",
						Description: "A tool for concurrent testing",
						MCPID:       "mcp-1",
						VerticalIDs: []string{"test-vertical"},
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		case "/mcp/v1/tools/concurrent-tool/execute":
			response := kirha.ExecuteToolResponse{
				Result: map[string]interface{}{
					"status":    "completed",
					"timestamp": time.Now().Unix(),
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &kirha.Config{
		APIKey:     "test-api-key",
		VerticalID: "test-vertical",
		BaseURL:    server.URL,
		Timeout:    time.Second * 5,
	}

	client := kirha.New(config)
	app := toolapp.New(client)

	ctx := context.Background()

	// Run multiple operations concurrently
	concurrency := 5
	errChan := make(chan error, concurrency*2)

	// List tools concurrently
	for i := 0; i < concurrency; i++ {
		go func() {
			_, err := app.ListTools(ctx)
			errChan <- err
		}()
	}

	// Execute tools concurrently
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			_, err := app.ExecuteTool(ctx, "concurrent-tool", mcpmux.Args{
				"request_id": id,
			})
			errChan <- err
		}(i)
	}

	// Check all operations completed successfully
	for i := 0; i < concurrency*2; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("concurrent operation failed: %v", err)
		}
	}
}
