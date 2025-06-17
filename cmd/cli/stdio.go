package cli

import (
	"context"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.kirha.ai/kirha-mcp-gateway/di"
)

// NewCmdStdio creates the stdio server command for the Kirha MCP Gateway.
// This command starts the gateway using standard input/output for MCP protocol
// communication. This mode is typically used when the gateway is invoked by
// another process (e.g., an IDE or LLM client) that communicates via pipes.
//
// The server configuration is loaded from environment variables:
//   - KIRHA_API_KEY: Required API key for Kirha API authentication
//   - KIRHA_VERTICAL: Required vertical ID for tool filtering
//   - KIRHA_BASE_URL: Kirha API base URL (default: https://api.kirha.ai)
//   - KIRHA_TIMEOUT: Request timeout (default: 120s)
//   - ENABLE_LOGS: Enable/disable logging (default: true)
//
// Returns:
//   - *cobra.Command: The stdio server command
//
// Example:
//
//	kirha-mcp-gateway stdio
func NewCmdStdio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start the MCP server",
		Long: `Start the Kirha MCP Gateway server that listens for MCP protocol
connections and proxies requests to the Kirha AI API.`,
		Example: "kirha-mcp-gateway stdio",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			l := logger.New("cli_stdio_mcp")

			l.Info("Starting stdio MCP server...")

			server, err := di.ProvideStdioServer()
			if err != nil {
				l.Error("Failed to initialize server", slog.String("error", err.Error()))
				return err
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
				<-sigCh
				l.Info("Received shutdown signal")

				shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := server.Stop(shutdownCtx); err != nil {
					l.Error("Failed to shutdown server gracefully", slog.String("error", err.Error()))
				}
				l.Info("Server shutdown complete")
				wg.Done()
			}()

			if err := server.Start(ctx); err != nil {
				l.Error("Server error", slog.String("error", err.Error()))
				return err
			}

			wg.Wait()
			l.Info("Server stopped")
			return nil
		},
	}
	return cmd
}
