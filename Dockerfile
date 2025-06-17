# Build stage
FROM golang:1.24-alpine AS builder

# Set up build environment
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install Wire
RUN go install github.com/google/wire/cmd/wire@latest

# Copy source code
COPY . .

# Generate Wire dependency injection code
RUN cd di && wire

# Build args for version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown
ARG GO_VERSION=unknown

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w \
      -X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.version=${VERSION}' \
      -X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.commit=${COMMIT}' \
      -X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.date=${DATE}' \
      -X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.goVersion=${GO_VERSION}'" \
    -o kirha-mcp-gateway ./cmd

# Final stage
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/kirha-mcp-gateway /usr/local/bin/kirha-mcp-gateway

# Create non-root user
USER 65534:65534

# Expose port (if needed for health checks)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/kirha-mcp-gateway"]

# Default command
CMD ["stdio"]

# Labels for metadata
LABEL org.opencontainers.image.title="Kirha MCP Gateway"
LABEL org.opencontainers.image.description="A high-performance Model Context Protocol server for Kirha AI"
LABEL org.opencontainers.image.url="https://go.kirha.ai/kirha-mcp-gateway"
LABEL org.opencontainers.image.source="https://go.kirha.ai/kirha-mcp-gateway"
LABEL org.opencontainers.image.vendor="Kirha AI"
LABEL org.opencontainers.image.licenses="MIT"