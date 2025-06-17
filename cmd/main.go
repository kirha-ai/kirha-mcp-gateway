package main

import (
	"errors"
	"github.com/joho/godotenv"
	"go.kirha.ai/kirha-mcp-gateway/cmd/cli"
	"go.kirha.ai/kirha-mcp-gateway/pkg/logger"
	"os"
)

// main is the application entry point.
// It loads environment variables from .env file (if present) and executes
// the root CLI command based on provided arguments.
func main() {
	log := logger.New("main")
	log.Info("Starting Kirha MCP Gateway")

	err := godotenv.Load(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		os.Exit(1)
	}

	rootCmd := cli.NewCmdRoot()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
