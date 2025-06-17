package cli

import (
	"context"
	"errors"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.kirha.ai/kirha-mcp-gateway/di"
)

// NewCmdHTTP creates the HTTP server command for the Kirha MCP Gateway.
// This command starts the gateway as an HTTP server that listens for MCP protocol
// requests over HTTP transport. The server implements graceful shutdown on SIGINT/SIGTERM.
//
// The server configuration is loaded from environment variables:
//   - SERVER_PORT: HTTP server port (default: 8080)
//   - KIRHA_API_KEY: Required API key for Kirha API authentication
//   - KIRHA_VERTICAL: Required vertical ID for tool filtering
//   - KIRHA_BASE_URL: Kirha API base URL (default: https://api.kirha.ai)
//   - KIRHA_TIMEOUT: Request timeout (default: 120s)
//
// Returns:
//   - *cobra.Command: The HTTP server command
//
// Example:
//
//	kirha-mcp-gateway http
func NewCmdHTTP() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Start the MCP server",
		Long: `Start the Kirha MCP Gateway server that listens for MCP protocol
connections and proxies requests to the Kirha AI API.`,
		Example: "kirha-mcp-gateway http",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			l := logger.New("cli_http_mcp")

			l.Info("Starting HTTP MCP server...")

			server, err := di.ProvideHTTPServer()
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

			if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
