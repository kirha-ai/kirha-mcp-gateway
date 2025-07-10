<p align="center">
  <a href="https://kirha.ai" target="_blank">
    <img src="assets/logo.png" width="200" alt="Kirha Logo">
  </a>
</p>

<h1 align="center">Kirha MCP Gateway</h1>
<h3 align="center">The AI bridge to reality</h3>


The golden light of private data will never shine on your AI queries if you can't find what you need and pay for it instantly.
**Kirha MCP Gateway** is an MCP (Model Context Protocol) that handles exactly that for you: routing and micropayment.

## Beta Scope

Kirha launched on the Crypto vertical with providers like Dune, Defillama, [Zerion](https://zerion.io/blog/how-kirha-leverages-zerion-api-to-revolutionize-crypto-data-access/), Xverse, Coingecko, Cielo...
Access live market intelligence & wallets' transactions and relations, whales' movements, sentiment analysis, and much more from the comfort of your favourite [AI Client](https://github.com/kirha-ai/mcp-installer?tab=readme-ov-file#supported-clients)
Kirha is under free Beta. Request an access code [here](https://app.kirha.ai/auth/claim-invite-code)

## Why Kirha MCP Gateway?

Kirha MCP Gateway is designed to simplify access to high-value, aggregated data sources through a single, unified interface.

- **One Auth to rule them all**: Access multiple premium data sources through a single MCP endpoint. You don't need to integrate or authenticate separately with multiple APIs.

- **Multi-API Composition**: In a single request, the gateway can compose and aggregate data from several APIs simultaneously. This reduces complexity and response time for complex queries.

- **Enhanced Context for Agents and Chatbots**: By centralizing and combining data, the gateway provides richer context for AI agents and conversational systems (such as Claude), improving their ability to deliver accurate, relevant responses.

## Tools

By default, the gateway operates in **auto mode**, where both planning and execution are performed automatically in a single request.
Planning is deterministic: semantically identical prompts yield the same composition of tool(s).
If the environment variable `TOOL_PLAN_MODE_ENABLED` is set to `true`, the gateway switches to **planning mode**, which separates the planning and execution steps.
When we ship Kirha's payment, you will be able to accept or reject the plan before it is executed.

### Auto mode

- `execute-tool-planning`: Automatically performs both planning and execution in a single request. This tool internally creates the plan and immediately executes it, returning the final result in one step.

### Plan mode

- `create-tool-planning`: Prepares a query and generates a plan_id for later execution. This tool allows creating a plan in advance, which can then be reviewed before being executed.
- `execute-tool-planning`: Executes a previously created plan by providing its plan_id. This step runs the predefined plan and returns the result.

## Installation

### Kirha MCP Installer (recommended)  

You can install and run the Kirha MCP Gateway using the [Kirha MCP Installer](https://github.com/kirha-ai/mcp-installer?tab=readme-ov-file#installation).

```
npx @kirha/mcp-installer install --client <client> --key <api-key>
```
