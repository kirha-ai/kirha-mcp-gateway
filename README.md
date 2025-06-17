# Kirha MCP Gateway

[![Go Version](https://img.shields.io/github/go-mod/go-version/kirha-ai/kirha-mcp-gateway)](https://golang.org)
[![License](https://img.shields.io/github/license/kirha-ai/kirha-mcp-gateway)](LICENSE)
[![Tests](https://go.kirha.ai/kirha-mcp-gateway/workflows/test/badge.svg)](https://go.kirha.ai/kirha-mcp-gateway/actions)

A Model Context Protocol (MCP) server that provides seamless access to Kirha AI tools. This gateway acts as a bridge between MCP clients and the Kirha API, enabling developers to integrate premium data providers into their applications with ease.

## Overview
The MCP Gateway is the entrypoint to all data providers featured with Kirha. They are grouped by verticals, which are specific domains of knowledge (e.g., crypto, finance, insurance, etc.). Providing a one-to-all solution for accessing multiple private data sources.

## Why?
LLM interfaces are not connected to private data in real time. The MCP protocol paves the way for LLMs to access private data sources in a standardized manner, but it requires to handle multiple API keys, endpoints, and subscriptions for each data provider.
With Kirha MCP Gateway any LLM with MCP support can access a wide range of data providers, without authenticating to each provider individually and pay only for what they use.
We run a cluster of MCP servers that expose data from our partner providers. Auth once to rule them all. Stop juggling multiple API keys, endpoints and subscriptions.
When the Gateway is used in external LLM interfaces (such as Claude, ChatGPT, or Cursor), we cannot apply smart routing as we do with our own interfaces. Read more about this in the [limitations](#limitations) section.

## 📦 Installation

### NPM (Recommended)

The npm package automatically downloads the appropriate binary for your platform:

```bash
# Run directly with npx (downloads binary on first use)
npx @kirha/mcp-gateway stdio

# Or install globally
npm install -g @kirha/mcp-gateway
kirha-mcp-gateway stdio
```

**Supported Platforms:**
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)  
- Windows (AMD64, ARM64)

### Download Binary

Download the latest binary for your platform from the [releases page](https://go.kirha.ai/kirha-mcp-gateway/releases).

### Build from Source

```bash
git clone https://go.kirha.ai/kirha-mcp-gateway.git
cd kirha-mcp-gateway
go build -o bin/kirha-mcp-gateway ./cmd
```

## 🚀 Quick Start

1. **Get your API key**
   - Get access to the Kirha, if you don't have an account you can request early access at [Kirha](https://app.kirha.ai).
   - Go through the [Kirha Documentation](https://kirha.gitbook.io/kirha-api) to obtain your API key.
   - Choose a vertical that suits your needs, available verticals are listed in the [Kirha documentation](https://kirha.gitbook.io/kirha-api/verticals).

2. **Set up your environment variables:**

```bash
export KIRHA_API_KEY="your-api-key"
export KIRHA_VERTICAL="your-vertical"
```

3. **Start the MCP server:**

```bash
# For stdio transport (typical MCP usage)
npx @kirha/mcp-gateway stdio

# For HTTP transport
npx @kirha/mcp-gateway http
```

4. **Connect your MCP client** to the gateway and start using Kirha AI tools!

## 🔌 MCP Client Example Configuration

### Claude Desktop

Add this to your Claude Desktop configuration file:

```json
{
  "mcpServers": {
    "kirha": {
      "command": "npx",
      "args": ["@kirha/mcp-gateway", "stdio"],
      "env": {
        "KIRHA_API_KEY": "your-api-key",
        "KIRHA_VERTICAL": "your-vertical"
      }
    }
  }
}
```

### Other MCP Clients

For other MCP clients, configure them to run:
```bash
kirha-mcp-gateway stdio
```

With the required environment variables set.

## ⚙️ Configuration

The gateway supports configuration through environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `KIRHA_API_KEY` | ✅ | - | Your Kirha AI API key |
| `KIRHA_VERTICAL` | ✅ | - | Your vertical for tool access |
| `KIRHA_TIMEOUT` | ❌ | `120s` | Request timeout duration |
| `ENABLE_LOGS` | ❌ | `true` | Enable/disable logging |
| `MCP_PORT` | ❌ | `8022` | HTTP server port (for HTTP transport) |
| `MCP_TOOL_CALL_TIMEOUT_SECONDS` | ❌ | `120` | Tool execution timeout in seconds |

### Running Locally

```bash
# Set environment variables
cp .env.example .env
# fill in your API key and vertical ID in .env

go run cmd/main.go stdio # or go run cmd/main.go http
```

## 📚 API Documentation

For detailed API documentation, refer to the [API Documentation](docs/API.md).

### MCP Protocol Support

The gateway implements the full [MCP 2025-03-26](https://modelcontextprotocol.io/specification/2025-03-26) specification:

- **Tools**: List and execute available Kirha AI tools
- **Resources**: Access to Kirha resources (planned)
- **Prompts**: Pre-configured prompts (planned)

### Available Tools

Tools are dynamically loaded based on your vertical configuration. Use the `list-tools` capability to see available tools for your account.

## Limitations
### Routing
Smart Kirha Routing is not supported in external LLM interfaces (e.g., Claude, ChatGPT, Cursor). The LLM will decide which tool to call based on the context, and the MCP server will execute the tool without any routing logic.
It will not be able to compose more than three tools in a single call. For smart routing, you need to use Kirha interface:
- API for agents: [Kirha API](https://kirha.gitbook.io/kirha-api/completion-api/chat-completion-openai)
- Kirha Chat for power users: [Kirha Chat](https://app.kirha.ai)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ by the [Kirha AI](https://kirha.ai) team