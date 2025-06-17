package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root command for the Kirha MCP Gateway CLI.
// It provides subcommands for different transport modes (HTTP and stdio).
//
// Returns:
//   - *cobra.Command: The root command with configured subcommands
//
// Available subcommands:
//   - http: Run the gateway as an HTTP server
//   - stdio: Run the gateway using standard input/output
func NewCmdRoot() *cobra.Command {
	var versionFlag bool

	cmd := &cobra.Command{
		Use:   "kirha-mcp-gateway",
		Short: "Kirha MCP Gateway - Connect to Kirha AI tools via MCP protocol",
		Long: `Kirha MCP Gateway is a Model Context Protocol (MCP) server that provides
access to Kirha AI tools and services. It acts as a bridge between MCP clients
and the Kirha AI API, allowing seamless integration of Kirha tools into
MCP-compatible applications.`,
		Version: Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if versionFlag {
				fmt.Printf("Kirha MCP Gateway version %s\n", Version)
				os.Exit(0)
			}
			
			// Check for updates on every command execution
			checkForUpdates()
		},
	}

	cmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "display version information")

	cmd.AddCommand(NewCmdHTTP())
	cmd.AddCommand(NewCmdStdio())
	cmd.AddCommand(NewCmdVersion())
	cmd.AddCommand(NewCmdUpdate())
	return cmd
}
