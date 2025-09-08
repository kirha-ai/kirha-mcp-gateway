#!/usr/bin/env node

import { serve } from "@hono/node-server";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

import {
  StreamableHTTPServerTransport,
  type StreamableHTTPServerTransportOptions,
} from "@modelcontextprotocol/sdk/server/streamableHttp.js";
import { toFetchResponse, toReqRes } from "fetch-to-node";
import { type Context, Hono } from "hono";
import { cors } from "hono/cors";
import { type Config, config } from "./config.js";
import { kirhaToolDefinitions } from "./kirha-tools.js";

function createStatelessServer(config: Config) {
  const server = new McpServer({
    name: config.mcpServer.name,
    version: config.mcpServer.version,
  });

  config.tools.forEach((tool) => {
    const toolDefinition = kirhaToolDefinitions[tool.name];

    server.registerTool(
      tool.name,
      {
        title: tool.name,
        description: tool.description,
        inputSchema: toolDefinition.inputSchema,
      },
      (input, extra) => {
        const apiKey = (extra.requestInfo?.headers["x-kirha-api-key"] as string) ?? config.apiKey;

        if (!apiKey) {
          throw new Error("KIRHA_API_KEY is required");
        }

        return toolDefinition.handler(input, { ...config, apiKey });
      },
    );
  });

  return server;
}

function streamableHTTPTransport(c: Context, options?: StreamableHTTPServerTransportOptions) {
  const transport = new StreamableHTTPServerTransport(
    options ?? {
      sessionIdGenerator: undefined,
    },
  );

  return Object.assign(transport, {
    async stream() {
      const { req, res } = toReqRes(c.req.raw);
      await transport.handleRequest(req, res, await c.req.json());
      return toFetchResponse(res);
    },
  });
}

function startHttpServer(config: Config) {
  const app = new Hono();
  const statelessMcpServer = createStatelessServer(config);

  app.use(
    cors({
      origin: "*",
      allowMethods: ["GET", "POST", "OPTIONS"],
      allowHeaders: ["Accept", "Content-Type", "x-kirha-api-key"],
      exposeHeaders: ["Content-Type"],
    }),
  );

  app.get("/health", (c) => {
    return c.json({ status: "healthy" });
  });

  app.post("/mcp", async (c) => {
    const transport = streamableHTTPTransport(c);
    await statelessMcpServer.connect(transport);
    return transport.stream();
  });

  serve(
    {
      fetch: app.fetch,
      port: config.port,
    },
    (info) => {
      const address = info.address === "::" ? "localhost" : `${info.address}`;
      const protocol = info.address === "::" ? "http" : "https";
      const url = `${protocol}://${address}:${info.port}`;

      console.log(`MCP server started on ${url}`);
    },
  );
}

async function startStdioServer(config: Config) {
  const server = createStatelessServer(config);
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

const isHttpMode = config.mcpServer.mode === "http";
const isStdioMode = config.mcpServer.mode === "stdio";

if (isHttpMode) {
  console.log("Start HTTP Mcp server");
  startHttpServer(config);
}

if (isStdioMode) {
  console.log("Start HTTP Mcp server");
  (async () => startStdioServer(config))();
}

process.on("SIGINT", async () => {
  console.log("Shutting down server...");
  process.exit(0);
});

process.on("SIGTERM", async () => {
  console.log("Shutting down server...");
  process.exit(0);
});
