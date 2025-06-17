package kirha

import (
	"context"
	"encoding/json"
	"errors"
	domainErrors "go.kirha.ai/kirha-mcp-gateway/internal/core/domain/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.kirha.ai/mcpmux"
)

func TestKirhaClient_ListTools(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse interface{}
		statusCode     int
		expectedError  error
		expectedCount  int
	}{
		{
			name: "successful response",
			serverResponse: ListToolsResponse{
				Tools: []ToolResponse{
					{
						ID:          "1",
						Identifier:  "test-tool",
						Name:        "Test Tool",
						Description: "A test tool",
						MCPID:       "mcp-1",
						VerticalIDs: []string{"vertical-1"},
					},
				},
			},
			statusCode:    http.StatusOK,
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:           "unauthorized",
			serverResponse: map[string]string{"error": "unauthorized"},
			statusCode:     http.StatusUnauthorized,
			expectedError:  domainErrors.ErrUnauthorized,
			expectedCount:  0,
		},
		{
			name:           "server error",
			serverResponse: map[string]string{"error": "internal server error"},
			statusCode:     http.StatusInternalServerError,
			expectedError:  domainErrors.ErrInternalServer,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			config := &Config{
				APIKey:     "test-key",
				VerticalID: "test-vertical",
				BaseURL:    server.URL,
				Timeout:    time.Second * 5,
			}

			client := New(config)
			tools, err := client.ListTools(context.Background())

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(tools) != tt.expectedCount {
				t.Errorf("expected %d tools, got %d", tt.expectedCount, len(tools))
			}
		})
	}
}

func TestKirhaClient_ExecuteTool(t *testing.T) {
	tests := []struct {
		name            string
		toolName        string
		arguments       mcpmux.Args
		serverResponse  interface{}
		statusCode      int
		expectedError   error
		expectedSuccess bool
	}{
		{
			name:      "successful execution",
			toolName:  "test-tool",
			arguments: mcpmux.Args{"param": "value"},
			serverResponse: ExecuteToolResponse{
				Result: map[string]interface{}{"output": "success"},
			},
			statusCode:      http.StatusOK,
			expectedError:   nil,
			expectedSuccess: true,
		},
		{
			name:            "tool not found",
			toolName:        "nonexistent-tool",
			arguments:       mcpmux.Args{},
			serverResponse:  map[string]string{"error": "tool not found"},
			statusCode:      http.StatusNotFound,
			expectedError:   domainErrors.ErrToolNotFound,
			expectedSuccess: false,
		},
		{
			name:            "unauthorized",
			toolName:        "test-tool",
			arguments:       mcpmux.Args{},
			serverResponse:  map[string]string{"error": "unauthorized"},
			statusCode:      http.StatusUnauthorized,
			expectedError:   domainErrors.ErrUnauthorized,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			config := &Config{
				APIKey:     "test-key",
				VerticalID: "test-vertical",
				BaseURL:    server.URL,
				Timeout:    time.Second * 5,
			}

			client := New(config)
			result, err := client.ExecuteTool(context.Background(), tt.toolName, tt.arguments)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != nil && result.Success != tt.expectedSuccess {
				t.Errorf("expected success %v, got %v", tt.expectedSuccess, result.Success)
			}
		})
	}
}
