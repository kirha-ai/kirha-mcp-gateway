//go:build wireinject
// +build wireinject

package di

import (
	"fmt"
	"github.com/google/wire"
	"go.kirha.ai/kirha-mcp-gateway/internal/adapters/clients/kirha"
	"go.kirha.ai/kirha-mcp-gateway/internal/adapters/handlers/mcp"
	"go.kirha.ai/kirha-mcp-gateway/internal/applications/toolapp"
	"go.kirha.ai/kirha-mcp-gateway/internal/core/ports/clients"
	"go.kirha.ai/mcpmux"
	"os"
)

func ProvideKirhaClient() (clients.KirhaClient, error) {
	wire.Build(
		kirha.LoadConfig,
		kirha.New,
	)

	return nil, nil
}

func ProvideToolApplication(k clients.KirhaClient) *toolapp.Application {
	return toolapp.New(k)
}

func ProvideStdioTransport() mcpmux.Transport {
	return mcpmux.NewStdioTransport(&mcpmux.StdioTransportConfig{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func ProvideHTTPConfig() (*mcp.HTTPConfig, error) {
	return mcp.LoadHTTPConfig()
}

func ProvideHTTPTransport(cfg *mcp.HTTPConfig) mcpmux.Transport {
	return mcpmux.NewHTTPTransport(&mcpmux.HTTPTransportConfig{
		Addr:            fmt.Sprintf(":%d", cfg.Port),
		Path:            "/mcp",
		AllowedOrigins:  []string{"*"},
		ResponseTimeout: cfg.ToolTimeout,
	})
}

var applicationSet = wire.NewSet(
	ProvideKirhaClient,
	ProvideToolApplication,
)

var httpSet = wire.NewSet(
	applicationSet,
	ProvideHTTPConfig,
	ProvideHTTPTransport,
)

var stdioSet = wire.NewSet(
	applicationSet,
	ProvideStdioTransport,
)

func ProvideHTTPServer() (*mcp.Server, error) {
	wire.Build(
		httpSet,
		mcp.New,
	)

	return nil, nil
}

func ProvideStdioServer() (*mcp.Server, error) {
	wire.Build(
		stdioSet,
		mcp.New,
	)

	return nil, nil
}
