# Kirha MCP Gateway API Documentation

This document describes the API and usage patterns for the Kirha MCP Gateway.

## Overview

The Kirha MCP Gateway implements the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) specification, providing access to Kirha AI tools through a standardized interface.

## Transport Modes

### Stdio Transport (Recommended)

The stdio transport is the recommended mode for MCP clients:

```bash
kirha-mcp-gateway stdio
```

This mode communicates via standard input/output and is typically used when the gateway is invoked by MCP clients like IDEs or AI applications.

### HTTP Transport

The HTTP transport mode runs an HTTP server:

```bash
kirha-mcp-gateway http
```

This mode listens on the configured port (default: 8022) for HTTP requests containing MCP protocol messages.

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `KIRHA_API_KEY` | ✅ | - | Your Kirha AI API key |
| `KIRHA_VERTICAL` | ✅ | - | Your vertical ID for tool access |
| `KIRHA_BASE_URL` | ❌ | `https://api.kirha.ai` | Base URL for Kirha API |
| `KIRHA_TIMEOUT` | ❌ | `120s` | Request timeout duration |
| `ENABLE_LOGS` | ❌ | `true` | Enable/disable logging |
| `MCP_PORT` | ❌ | `8022` | HTTP server port (for HTTP transport) |
| `MCP_TOOL_CALL_TIMEOUT_SECONDS` | ❌ | `120` | Tool execution timeout in seconds |

### Example Configuration

```bash
export KIRHA_API_KEY="your-api-key-here"
export KIRHA_VERTICAL="your-vertical-id"
export KIRHA_BASE_URL="https://api.kirha.ai"
export KIRHA_TIMEOUT="60s"
export ENABLE_LOGS="true"
```

## MCP Protocol Support

The gateway implements the following MCP capabilities:

### Tools

- **`list_tools`**: Retrieve all available tools for your vertical
- **`call_tool`**: Execute a specific tool with provided arguments

### Future Capabilities (Planned)

- **Resources**: Access to Kirha resources
- **Prompts**: Pre-configured prompt templates

## Tool Discovery

### Listing Tools

Tools are dynamically discovered based on your vertical configuration. Use the MCP `list_tools` capability to see available tools:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
```

### Tool Schema

Each tool includes:

```json
{
  "name": "tool-identifier",
  "description": "Description of what the tool does",
  "inputSchema": {
    "type": "object",
    "properties": {
      "parameter1": {
        "type": "string",
        "description": "Description of parameter1"
      }
    },
    "required": ["parameter1"]
  }
}
```

## Tool Execution

### Executing Tools

Execute tools using the MCP `call_tool` capability:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "tool-identifier",
    "arguments": {
      "parameter1": "value1",
      "parameter2": "value2"
    }
  }
}
```

### Tool Results

Tool execution returns structured results:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Tool execution result data"
      }
    ]
  }
}
```

## Client Integration

### Claude Desktop

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key",
        "KIRHA_VERTICAL": "your-vertical-id"
      }
    }
  }
}
```

### Continue.dev

Add to your Continue configuration:

```json
{
  "mcpServers": [
    {
      "name": "kirha",
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key",
        "KIRHA_VERTICAL": "your-vertical-id"
      }
    }
  ]
}
```

### Custom MCP Clients

For custom MCP clients, run the gateway in stdio mode and communicate via JSON-RPC over stdin/stdout:

```bash
kirha-mcp-gateway stdio
```

## Error Handling

### Error Responses

The gateway returns standard MCP error responses:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "Tool not found",
    "data": {
      "type": "TOOL_NOT_FOUND"
    }
  }
}
```

### Common Error Codes

| Code | Type | Description |
|------|------|-------------|
| -32000 | TOOL_NOT_FOUND | Requested tool does not exist |
| -32001 | UNAUTHORIZED | Invalid API key or vertical |
| -32002 | TIMEOUT | Request timed out |
| -32003 | INVALID_ARGUMENTS | Tool arguments are invalid |
| -32603 | INTERNAL_ERROR | Internal server error |

## Logging

### Log Format

The gateway uses structured JSON logging compatible with Google Cloud Platform:

```json
{
  "time": "2024-01-15T10:30:45.123Z",
  "severity": "INFO",
  "message": "Tool executed successfully",
  "application": "kirha-mcp-gateway",
  "logger": "mcp_handler",
  "tool_name": "example-tool",
  "duration": "1.234s"
}
```

### Log Levels

- **DEBUG**: Detailed debugging information
- **INFO**: General operational messages
- **WARNING**: Warning conditions
- **ERROR**: Error conditions that need attention

### Disabling Logs

For production stdio usage where logs might interfere with MCP communication:

```bash
ENABLE_LOGS=false kirha-mcp-gateway stdio
```

## Health Checks

### HTTP Transport Health

For HTTP transport, the gateway provides basic health monitoring through successful MCP operations.

### Stdio Transport Health

For stdio transport, health is verified through successful tool listing operations.

## Performance Considerations

### Timeouts

- **Request Timeout**: Configurable via `KIRHA_TIMEOUT` (default: 120s)
- **Tool Execution Timeout**: Configurable via `MCP_TOOL_CALL_TIMEOUT_SECONDS` (default: 120s)

### Concurrency

The gateway supports concurrent tool execution with proper error isolation.

### Caching

Tool schemas are cached per session to improve performance for repeated `list_tools` calls.

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify `KIRHA_API_KEY` is set correctly
   - Ensure `KIRHA_VERTICAL` matches your account configuration

2. **Tool Not Found Errors**
   - Check that tools are available for your vertical
   - Verify tool names match exactly (case-sensitive)

3. **Timeout Errors**
   - Increase `KIRHA_TIMEOUT` for slow network connections
   - Check Kirha API service status

4. **MCP Client Connection Issues**
   - Ensure the gateway binary is in your PATH
   - Verify environment variables are properly set in client configuration
   - Check client logs for specific error messages

### Debug Mode

Enable debug logging for troubleshooting:

```bash
ENABLE_LOGS=true kirha-mcp-gateway stdio
```

## Examples

### Complete Integration Example

1. **Install the gateway:**
```bash
npm install -g @kirha/mcp-server
```

2. **Set up environment:**
```bash
export KIRHA_API_KEY="sk-..."
export KIRHA_VERTICAL="my-vertical"
```

3. **Configure Claude Desktop:**
```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "sk-...",
        "KIRHA_VERTICAL": "my-vertical"
      }
    }
  }
}
```

4. **Use tools in Claude:**
- Claude will automatically discover available Kirha tools
- Tools can be invoked through natural language
- Results are returned in a structured format

### Direct API Usage

For direct API testing or custom integrations:

```bash
# Start HTTP server
KIRHA_API_KEY=sk-... KIRHA_VERTICAL=my-vertical kirha-mcp-gateway http

# Test tool listing
curl -X POST http://localhost:8022 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list"
  }'
```

## API Reference

### Method: `tools/list`

List all available tools for the configured vertical.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "example-tool",
        "description": "An example tool",
        "inputSchema": {
          "type": "object",
          "properties": {
            "input": {
              "type": "string",
              "description": "Input parameter"
            }
          },
          "required": ["input"]
        }
      }
    ]
  }
}
```

### Method: `tools/call`

Execute a specific tool with provided arguments.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "example-tool",
    "arguments": {
      "input": "Hello, world!"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Tool executed successfully with result: Hello, world!"
      }
    ]
  }
}
```

## Security Best Practices

1. **API Key Management**
   - Store API keys securely (environment variables, not code)
   - Rotate API keys regularly
   - Use separate keys for development and production

2. **Network Security**
   - Use HTTPS for all API communications
   - Configure appropriate timeouts
   - Monitor for unusual usage patterns

3. **Input Validation**
   - Tool arguments are validated against schemas
   - Input sanitization is performed automatically
   - Error messages are sanitized to prevent information leakage

## Support

For API support and questions:

- 📧 Email: developers@kirha.ai
- 📖 Documentation: [https://kirha.gitbook.io/kirha-api](https://docs.kirha.ai)
- 🐛 Issues: [GitHub Issues](https://go.kirha.ai/kirha-mcp-gateway/issues)