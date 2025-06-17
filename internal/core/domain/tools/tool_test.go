package tools

import (
	"reflect"
	"testing"
	"time"
)

// TestTool tests the Tool struct and its fields.
func TestTool(t *testing.T) {
	tests := []struct {
		name string
		tool Tool
	}{
		{
			name: "complete tool",
			tool: Tool{
				ID:          "test-id-1",
				Identifier:  "test-tool",
				Name:        "Test Tool",
				Description: "A comprehensive test tool",
				MCPID:       "mcp-test-1",
				VerticalIDs: []string{"vertical-1", "vertical-2"},
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"input": map[string]interface{}{
							"type":        "string",
							"description": "Input parameter",
						},
					},
					"required": []string{"input"},
				},
				Outputs: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"result": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		{
			name: "minimal tool",
			tool: Tool{
				ID:          "minimal-id",
				Name:        "Minimal Tool",
				Description: "A minimal tool",
			},
		},
		{
			name: "tool with empty arrays",
			tool: Tool{
				ID:          "empty-arrays-id",
				Name:        "Empty Arrays Tool",
				Description: "Tool with empty arrays",
				VerticalIDs: []string{},
			},
		},
		{
			name: "tool with nil parameters",
			tool: Tool{
				ID:          "nil-params-id",
				Name:        "Nil Parameters Tool",
				Description: "Tool with nil parameters",
				Parameters:  nil,
				Outputs:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := tt.tool

			// Test field access
			if tool.ID == "" && tt.name != "minimal tool" {
				t.Error("expected non-empty ID")
			}

			if tool.Name == "" {
				t.Error("expected non-empty Name")
			}

			if tool.Description == "" {
				t.Error("expected non-empty Description")
			}

			// Test that fields can be modified
			originalName := tool.Name
			tool.Name = "Modified Name"
			if tool.Name != "Modified Name" {
				t.Error("failed to modify tool name")
			}
			tool.Name = originalName // Restore

			// Test that arrays are properly handled
			if tool.VerticalIDs != nil {
				// Should be able to append to the slice
				originalLen := len(tool.VerticalIDs)
				tool.VerticalIDs = append(tool.VerticalIDs, "new-vertical")
				if len(tool.VerticalIDs) != originalLen+1 {
					t.Error("failed to append to VerticalIDs")
				}
			}

			// Test that maps can be modified
			if tool.Parameters != nil {
				params, ok := tool.Parameters.(map[string]interface{})
				if ok {
					params["new_param"] = "new_value"
					if params["new_param"] != "new_value" {
						t.Error("failed to modify parameters map")
					}
				}
			}
		})
	}
}

// TestToolEquality tests comparing Tool structs.
func TestToolEquality(t *testing.T) {
	tool1 := Tool{
		ID:          "test-1",
		Identifier:  "test-tool",
		Name:        "Test Tool",
		Description: "A test tool",
		MCPID:       "mcp-1",
		VerticalIDs: []string{"vertical-1"},
	}

	tool2 := Tool{
		ID:          "test-1",
		Identifier:  "test-tool",
		Name:        "Test Tool",
		Description: "A test tool",
		MCPID:       "mcp-1",
		VerticalIDs: []string{"vertical-1"},
	}

	tool3 := Tool{
		ID:          "test-2",
		Identifier:  "different-tool",
		Name:        "Different Tool",
		Description: "A different tool",
		MCPID:       "mcp-2",
		VerticalIDs: []string{"vertical-2"},
	}

	// Test equality with reflect.DeepEqual
	if !reflect.DeepEqual(tool1, tool2) {
		t.Error("expected tool1 and tool2 to be equal")
	}

	if reflect.DeepEqual(tool1, tool3) {
		t.Error("expected tool1 and tool3 to be different")
	}

	// Test field-by-field comparison
	if tool1.ID != tool2.ID {
		t.Error("expected same ID")
	}

	if tool1.ID == tool3.ID {
		t.Error("expected different ID")
	}
}

// TestToolExecutionResult tests the ToolExecutionResult struct.
func TestToolExecutionResult(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name   string
		result ToolExecutionResult
	}{
		{
			name: "successful execution",
			result: ToolExecutionResult{
				ToolName: "test-tool",
				Result: map[string]interface{}{
					"output": "success",
					"data":   []string{"a", "b", "c"},
				},
				Duration:  150 * time.Millisecond,
				Success:   true,
				Error:     "",
				Timestamp: now,
			},
		},
		{
			name: "failed execution",
			result: ToolExecutionResult{
				ToolName:  "failing-tool",
				Result:    nil,
				Duration:  50 * time.Millisecond,
				Success:   false,
				Error:     "Tool execution failed: invalid input",
				Timestamp: now,
			},
		},
		{
			name: "execution with empty result",
			result: ToolExecutionResult{
				ToolName:  "empty-tool",
				Result:    map[string]interface{}{},
				Duration:  10 * time.Millisecond,
				Success:   true,
				Error:     "",
				Timestamp: now,
			},
		},
		{
			name: "execution with complex result",
			result: ToolExecutionResult{
				ToolName: "complex-tool",
				Result: map[string]interface{}{
					"metadata": map[string]interface{}{
						"version":   "1.0",
						"timestamp": now.Unix(),
						"status":    "completed",
					},
					"data": []interface{}{
						map[string]interface{}{"id": 1, "value": "first"},
						map[string]interface{}{"id": 2, "value": "second"},
					},
					"summary": "Processed 2 items successfully",
				},
				Duration:  300 * time.Millisecond,
				Success:   true,
				Error:     "",
				Timestamp: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result

			// Test basic field access
			if result.ToolName == "" {
				t.Error("expected non-empty ToolName")
			}

			if result.Duration < 0 {
				t.Error("expected non-negative Duration")
			}

			if result.Timestamp.IsZero() {
				t.Error("expected non-zero Timestamp")
			}

			// Test success/error consistency
			if result.Success && result.Error != "" {
				t.Error("successful execution should not have error message")
			}

			if !result.Success && result.Error == "" {
				t.Error("failed execution should have error message")
			}

			// Test result modification
			if result.Result != nil {
				result.Result["test_modification"] = "test_value"
				if result.Result["test_modification"] != "test_value" {
					t.Error("failed to modify result map")
				}
			}

			// Test immutability of timestamp
			originalTime := result.Timestamp
			result.Timestamp = time.Now().Add(time.Hour)
			if result.Timestamp.Equal(originalTime) {
				t.Error("timestamp should be modifiable")
			}
		})
	}
}

// TestToolExecutionResultEquality tests comparing ToolExecutionResult structs.
func TestToolExecutionResultEquality(t *testing.T) {
	timestamp := time.Now()
	
	result1 := ToolExecutionResult{
		ToolName: "test-tool",
		Result: map[string]interface{}{
			"output": "test",
		},
		Duration:  100 * time.Millisecond,
		Success:   true,
		Error:     "",
		Timestamp: timestamp,
	}

	result2 := ToolExecutionResult{
		ToolName: "test-tool",
		Result: map[string]interface{}{
			"output": "test",
		},
		Duration:  100 * time.Millisecond,
		Success:   true,
		Error:     "",
		Timestamp: timestamp,
	}

	result3 := ToolExecutionResult{
		ToolName: "different-tool",
		Result: map[string]interface{}{
			"output": "different",
		},
		Duration:  200 * time.Millisecond,
		Success:   false,
		Error:     "failed",
		Timestamp: timestamp.Add(time.Minute),
	}

	// Test equality
	if !reflect.DeepEqual(result1, result2) {
		t.Error("expected result1 and result2 to be equal")
	}

	if reflect.DeepEqual(result1, result3) {
		t.Error("expected result1 and result3 to be different")
	}

	// Test individual field comparisons
	if result1.Success != result2.Success {
		t.Error("expected same success status")
	}

	if result1.Success == result3.Success {
		t.Error("expected different success status")
	}
}