#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";
import { InMemoryEventStore } from "@modelcontextprotocol/sdk/examples/shared/inMemoryEventStore.js";
import { isInitializeRequest } from "@modelcontextprotocol/sdk/types.js";
import { randomUUID } from "node:crypto";
import express from "express";
import { z } from "zod";

import { registerKirhaTools } from "./tools/kirha.js";

// Configuration schema
const ConfigSchema = z.object({
  kirhaApiKey: z.string().min(1, "Kirha API key is required"),
  verticalId: z.string().default("crypto"),
  toolPlanModeEnabled: z.boolean().default(false),
});

export type Config = z.infer<typeof ConfigSchema>;

/**
 * Create and configure an MCP server with Kirha tools
 */
export function createServer(config: Config) {
  const server = new McpServer({
    name: "kirha-mcp-gateway",
    version: "0.0.6",
  });

  // Register Kirha tools
  registerKirhaTools(server, config);

  return server;
}

/**
 * Extract configuration from query parameters
 */
function extractConfig(query: any): Config {
  const rawConfig = {
    kirhaApiKey: query.kirhaApiKey,
    verticalId: query.verticalId || "crypto",
    toolPlanModeEnabled: query.toolPlanModeEnabled === "true" || query.toolPlanModeEnabled === true,
  };

  try {
    return ConfigSchema.parse(rawConfig);
  } catch (error) {
    console.error("Invalid configuration:", error);
    throw new Error("Invalid configuration provided");
  }
}

// Check if running in HTTP mode or STDIO mode
const isHttpMode = process.env.PORT || process.argv.includes('--http');

if (isHttpMode) {
  // HTTP mode for Smithery
  const app = express();
  app.use(express.json({ limit: '10mb' }));
  
  // Store transports and servers by session ID
  const transports: Record<string, StreamableHTTPServerTransport> = {};
  const servers: Record<string, McpServer> = {};
  
  // MCP endpoint for POST requests
  app.post('/mcp', async (req, res) => {
    try {
      let transport: StreamableHTTPServerTransport;
      const sessionId = req.headers['mcp-session-id'] as string;
      
      // Extract and validate configuration from query parameters
      const config = extractConfig(req.query);

      if (sessionId && transports[sessionId]) {
        // Reuse existing transport
        transport = transports[sessionId];
      } else if (!sessionId && isInitializeRequest(req.body)) {
        // New initialization request
        const eventStore = new InMemoryEventStore();
        transport = new StreamableHTTPServerTransport({
          sessionIdGenerator: () => randomUUID(),
          eventStore,
          onsessioninitialized: (sessionId: string) => {
            console.log(`Session initialized with ID: ${sessionId}`);
            transports[sessionId] = transport;
          }
        });
        
        // Set up onclose handler
        transport.onclose = () => {
          const sid = transport.sessionId;
          if (sid) {
            console.log(`Transport closed for session ${sid}`);
            delete transports[sid];
            delete servers[sid];
          }
        };
        
        // Create server with configuration
        const server = createServer(config);
        
        // Store the server instance
        servers[transport.sessionId!] = server;
        
        await server.connect(transport);
        await transport.handleRequest(req, res, req.body);
        return;
      } else {
        res.status(400).json({
          jsonrpc: '2.0',
          error: {
            code: -32000,
            message: 'Bad Request: No valid session ID provided',
          },
          id: null,
        });
        return;
      }
      
      await transport.handleRequest(req, res, req.body);
    } catch (error) {
      console.error('Error handling MCP request:', error);
      if (!res.headersSent) {
        res.status(500).json({
          jsonrpc: '2.0',
          error: {
            code: -32603,
            message: error instanceof Error ? error.message : 'Internal server error',
          },
          id: null,
        });
      }
    }
  });
  
  // Handle GET requests for SSE streams
  app.get('/mcp', async (req, res) => {
    const sessionId = req.headers['mcp-session-id'] as string;
    
    if (!sessionId || !transports[sessionId]) {
      res.status(400).send('Invalid or missing session ID');
      return;
    }
    
    const transport = transports[sessionId];
    
    // Check for Last-Event-ID header for resumability
    const lastEventId = req.headers['last-event-id'];
    if (lastEventId) {
      console.log(`Client reconnecting with Last-Event-ID: ${lastEventId}`);
    }
    
    // Handle the SSE request
    await transport.handleRequest(req, res);
  });
  
  // Health check endpoint
  app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
  });
  
  // Start server
  const port = parseInt(process.env.PORT || "3000", 10);
  app.listen(port, '0.0.0.0', () => {
    console.log(`MCP HTTP server listening on port ${port}`);
  });
} else {
  // STDIO mode for Claude Desktop
  const kirhaApiKey = process.env.KIRHA_API_KEY;
  const verticalId = process.env.VERTICAL_ID || "crypto";
  const toolPlanModeEnabled = process.env.TOOL_PLAN_MODE_ENABLED === "true";
  
  if (!kirhaApiKey) {
    console.error("KIRHA_API_KEY environment variable is required");
    process.exit(1);
  }

  const config: Config = {
    kirhaApiKey,
    verticalId,
    toolPlanModeEnabled,
  };

  const server = createServer(config);
  const transport = new StdioServerTransport();
  server.connect(transport);
}