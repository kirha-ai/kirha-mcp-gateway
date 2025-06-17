# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-06-17

### Added
- Initial public release of Kirha MCP Gateway
- Model Context Protocol (MCP) server implementation with stdio and HTTP transports
- Hexagonal architecture with clean separation of concerns
- HTTP client for Kirha AI API integration with retry logic
- Comprehensive error handling and structured logging (GCP-compatible JSON format)
- CLI interface with Cobra framework (`stdio` and `http` commands)
- Multi-platform binary builds (Linux, macOS, Windows, ARM64)
- NPM package distribution for easy installation
- Docker support with multi-stage builds
- Comprehensive test coverage (unit, integration, and end-to-end tests)
- Tool discovery and execution through MCP protocol
- Configuration via environment variables with validation
- Wire-based dependency injection for clean architecture
- GitHub Actions CI/CD pipeline with automated testing
- Comprehensive documentation with examples and API references
- MCP client configuration examples (Claude Desktop, etc.)

### Features
- **Tool Management**: List and execute Kirha AI tools through MCP protocol
- **Dual Transport**: Support for both stdio and HTTP transport modes
- **Dynamic Tool Loading**: Tools are loaded based on vertical configuration
- **Request Tracing**: Structured logging with request IDs and context
- **Graceful Shutdown**: Proper cleanup on SIGINT/SIGTERM signals
- **Concurrent Processing**: Safe concurrent tool execution
- **Error Mapping**: Domain errors properly mapped to MCP error responses

### Developer Experience
- Comprehensive godoc documentation for all public APIs
- Table-driven tests with descriptive scenarios
- Mock implementations for testing external dependencies
- Integration tests covering full request/response cycles
- Development setup with clear contribution guidelines
- Code examples and usage patterns in documentation

### Security
- Secure API key handling with environment variable configuration
- Request/response validation and sanitization
- Timeout protection for HTTP requests (configurable)
- Input validation for tool arguments
- Error message sanitization to prevent information leakage