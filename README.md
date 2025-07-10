<p align="center">
  <a href="https://kirha.ai" target="_blank">
    <img src="assets/logo.png" width="200" alt="Bright Data Logo">
  </a>
</p>

<h1 align="center">Kirha MCP Gateway</h1>
<h3 align="center">The AI bridge to reality</h3>


**Kirha MCP Gateway** is an MCP (Model Context Protocol) server that provides access to **Kirha's Tool Planning** system.

## Why Kirha MCP Gateway?

Kirha MCP Gateway is designed to simplify access to high-value, aggregated data sources through a single, unified interface. Here are the key advantages:

- **Unified Access Point**: Access multiple premium data sources through a single MCP endpoint. No need to integrate or authenticate separately with multiple APIs.

- **Multi-API Composition**: In a single request, the gateway can compose and aggregate data from several APIs simultaneously. This reduces complexity and response time for complex queries.

- **Enhanced Context for Agents and Chatbots**: By centralizing and combining data, the gateway provides richer context for AI agents and conversational systems (such as Claude), improving their ability to deliver accurate, relevant responses.

## Tools

By default, the gateway operates in **auto mode**, where both planning and execution are performed automatically in a single request.
If the environment variable `TOOL_PLAN_MODE_ENABLED` is set to `true`, the gateway switches to **planning mode**, which separates the planning and execution steps.

### Auto mode

- `execute-tool-planning`: Automatically performs both planning and execution in a single request. This tool internally creates the plan and immediately executes it, returning the final result in one step.

### Plan mode

- `create-tool-planning`: Prepares a query and generates a plan_id for later execution. This tool allows creating a plan in advance, which can then be reviewed before being executed.
- `execute-tool-planning`: Executes a previously created plan by providing its plan_id. This step runs the predefined plan and returns the result.

## Installation

### Kirha Mcp Installer (recommended)  

You can install and run the Kirha MCP Gateway using the [Kirha MCP Installer](https://github.com/kirha-ai/mcp-installer?tab=readme-ov-file#installation).

```
npx @kirha/mcp-installer install --client <client> --key <api-key>
```
