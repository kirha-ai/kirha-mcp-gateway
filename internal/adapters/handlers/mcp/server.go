package mcp

import (
	"context"
	"go.kirha.ai/kirha-mcp-gateway/internal/adapters/handlers/mcp/toolshdl"
	"go.kirha.ai/kirha-mcp-gateway/internal/applications/toolapp"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"go.kirha.ai/mcpmux"
	"log/slog"
)

// Server implements an MCP server that exposes Kirha tools through the MCP protocol.
// It handles tool discovery and execution requests from MCP clients.
type Server struct {
	toolApp   *toolapp.Application
	logger    *slog.Logger
	server    *mcpmux.Server
	transport mcpmux.Transport
}

// New creates a new MCP server instance with the provided tool application and transport.
// The server will use the tool application to handle tool operations and the transport
// for MCP protocol communication.
//
// Parameters:
//   - toolApp: The application layer for tool management operations
//   - transport: The transport layer for MCP communication (HTTP or stdio)
//
// Returns:
//   - *Server: A new MCP server instance ready to be started
func New(toolApp *toolapp.Application, transport mcpmux.Transport) *Server {
	return &Server{
		toolApp:   toolApp,
		logger:    logger.New("mcp_server"),
		transport: transport,
	}
}

// Start initializes and starts the MCP server with the configured transport.
// It creates a new mcpmux server with appropriate handlers and begins serving requests.
// The server will continue running until an error occurs or Stop is called.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation
//
// Returns:
//   - error: An error if the server fails to start or encounters an error during operation
func (s *Server) Start(ctx context.Context) error {
	mcpHdl := toolshdl.New(s.toolApp)
	server := mcpmux.NewServer(
		mcpmux.WithServerInfo("Kirha MCP", "0.1.0"),
		mcpmux.WithInstructions("Gateway to premium data providers for real time insights"),
		mcpmux.WithLogger(s.logger),
	)
	server.UseToolAggregator(mcpHdl)

	s.server = server
	s.logger.InfoContext(ctx, "Starting mcp server")
	err := s.server.Serve(s.transport)
	return err
}

// Stop gracefully shuts down the MCP server.
// It ensures all active connections are closed and resources are released.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation
//
// Returns:
//   - error: An error if the server fails to stop gracefully
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	if err := s.server.Stop(); err != nil {
		s.logger.Error("failed to stop MCP server", slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("MCP server stopped successfully")
	return nil
}
