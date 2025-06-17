# Claude Desktop Integration Example

This example shows how to integrate the Kirha MCP Gateway with Claude Desktop.

## Prerequisites

1. Install Claude Desktop from [Claude.ai](https://claude.ai/download)
2. Install the Kirha MCP Gateway:
   ```bash
   npm install -g @kirha/mcp-server
   ```
3. Get your Kirha API credentials (API key and vertical ID)

## Configuration

### 1. Locate Configuration File

The Claude Desktop configuration file location depends on your operating system:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

### 2. Add Kirha MCP Server

Edit the configuration file to include the Kirha MCP server:

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key-here",
        "KIRHA_VERTICAL": "your-vertical-id"
      }
    }
  }
}
```

### 3. Complete Configuration Example

Here's a complete configuration file with Kirha and optional additional settings:

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "sk-1234567890abcdef...",
        "KIRHA_VERTICAL": "finance",
        "KIRHA_BASE_URL": "https://api.kirha.ai",
        "KIRHA_TIMEOUT": "120s",
        "ENABLE_LOGS": "false"
      }
    }
  }
}
```

## Usage

### 1. Restart Claude Desktop

After saving the configuration, restart Claude Desktop to load the new MCP server.

### 2. Verify Integration

You should see Kirha tools available in Claude Desktop. You can verify this by:

1. Starting a new conversation
2. Asking Claude: "What Kirha tools are available?"
3. Claude should list the tools from your configured vertical

### 3. Using Tools

Once integrated, you can use Kirha tools naturally in conversation:

```
You: "Use the market data tool to get the latest price for AAPL"

Claude: I'll use the Kirha market data tool to get the latest price for Apple Inc. (AAPL).

[Claude automatically calls the appropriate Kirha tool]

The latest price for AAPL is $175.43, up 2.3% from the previous close...
```

## Advanced Configuration

### Custom Tool Timeout

For tools that might take longer to execute:

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key-here",
        "KIRHA_VERTICAL": "your-vertical-id",
        "MCP_TOOL_CALL_TIMEOUT_SECONDS": "300"
      }
    }
  }
}
```

### Multiple Verticals

If you have access to multiple verticals, you can configure separate servers:

```json
{
  "mcpServers": {
    "kirha-finance": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key-here",
        "KIRHA_VERTICAL": "finance"
      }
    },
    "kirha-healthcare": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key-here",
        "KIRHA_VERTICAL": "healthcare"
      }
    }
  }
}
```

## Troubleshooting

### Common Issues

1. **Tools not appearing in Claude Desktop**
   - Check that the configuration file syntax is valid JSON
   - Verify the file path is correct for your operating system
   - Restart Claude Desktop after making changes
   - Check that `@kirha/mcp-server` is installed globally

2. **Authentication errors**
   - Verify your `KIRHA_API_KEY` is correct
   - Ensure your `KIRHA_VERTICAL` matches your account setup
   - Check that your API key has access to the specified vertical

3. **Tool execution timeouts**
   - Increase `MCP_TOOL_CALL_TIMEOUT_SECONDS` for slower tools
   - Check your internet connection
   - Verify Kirha API service status

### Debug Mode

Enable logging to troubleshoot issues:

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-server", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key-here",
        "KIRHA_VERTICAL": "your-vertical-id",
        "ENABLE_LOGS": "true"
      }
    }
  }
}
```

Note: Logs will be written to Claude Desktop's console. You can view them by opening Developer Tools in Claude Desktop.

### Validation

Test your configuration manually:

```bash
# Set environment variables
export KIRHA_API_KEY="your-api-key-here"
export KIRHA_VERTICAL="your-vertical-id"

# Test the gateway directly
npx @kirha/mcp-server stdio
```

The gateway should start without errors and be ready to receive MCP commands.

## Example Conversations

### Using Financial Tools

```
You: "What's the current market cap of Microsoft?"

Claude: I'll use the Kirha financial data tool to get Microsoft's current market capitalization.

[Tool execution result]
Microsoft (MSFT) currently has a market capitalization of $2.89 trillion based on the latest trading data...
```

### Using Data Analysis Tools

```
You: "Analyze the trend in Tesla's stock price over the last month"

Claude: I'll use Kirha's market analysis tools to examine Tesla's stock price trend over the past month.

[Tool execution result]
Based on the analysis of Tesla's stock price over the last 30 days:
- Starting price: $185.43
- Ending price: $201.76
- Change: +8.8%
- Volatility: 2.3% daily average...
```

## Best Practices

1. **Security**
   - Never commit API keys to version control
   - Use environment variables or secure credential storage
   - Regularly rotate API keys

2. **Performance**
   - Set appropriate timeouts for your use case
   - Monitor tool usage to avoid rate limits
   - Use specific tool names when possible

3. **Organization**
   - Use descriptive server names for multiple verticals
   - Group related tools by vertical
   - Document your configuration for team sharing

## Getting Help

If you encounter issues:

1. Check the [API Documentation](../API.md)
2. Review the [troubleshooting guide](../troubleshooting.md)
3. Contact support at support@kirha.ai
4. Open an issue on [GitHub](https://go.kirha.ai/kirha-mcp-gateway/issues)