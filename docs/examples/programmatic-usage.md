# Programmatic Usage Examples

This document shows how to integrate the Kirha MCP Gateway programmatically into your applications.

## Node.js Integration

### Basic Setup

```javascript
import { spawn } from 'child_process';

class KirhaMCPClient {
  constructor(apiKey, verticalId) {
    this.apiKey = apiKey;
    this.verticalId = verticalId;
    this.process = null;
    this.requestId = 1;
  }

  async start() {
    // Start the Kirha MCP Gateway in stdio mode
    this.process = spawn('npx', ['@kirha/mcp-server', 'stdio'], {
      env: {
        ...process.env,
        KIRHA_API_KEY: this.apiKey,
        KIRHA_VERTICAL: this.verticalId,
        ENABLE_LOGS: 'false' // Disable logs for stdio mode
      },
      stdio: ['pipe', 'pipe', 'pipe']
    });

    // Set up JSON-RPC communication
    this.setupCommunication();
  }

  setupCommunication() {
    let buffer = '';
    
    this.process.stdout.on('data', (data) => {
      buffer += data.toString();
      
      // Process complete JSON messages
      const lines = buffer.split('\n');
      buffer = lines.pop(); // Keep incomplete line in buffer
      
      lines.forEach(line => {
        if (line.trim()) {
          try {
            const response = JSON.parse(line);
            this.handleResponse(response);
          } catch (err) {
            console.error('Failed to parse JSON:', err, line);
          }
        }
      });
    });

    this.process.stderr.on('data', (data) => {
      console.error('Gateway error:', data.toString());
    });
  }

  async sendRequest(method, params = {}) {
    return new Promise((resolve, reject) => {
      const id = this.requestId++;
      const request = {
        jsonrpc: '2.0',
        id,
        method,
        params
      };

      // Store resolver for this request
      this.pendingRequests = this.pendingRequests || new Map();
      this.pendingRequests.set(id, { resolve, reject });

      // Send request
      this.process.stdin.write(JSON.stringify(request) + '\n');
    });
  }

  handleResponse(response) {
    if (response.id && this.pendingRequests?.has(response.id)) {
      const { resolve, reject } = this.pendingRequests.get(response.id);
      this.pendingRequests.delete(response.id);

      if (response.error) {
        reject(new Error(response.error.message));
      } else {
        resolve(response.result);
      }
    }
  }

  async listTools() {
    return this.sendRequest('tools/list');
  }

  async callTool(name, arguments) {
    return this.sendRequest('tools/call', { name, arguments });
  }

  async stop() {
    if (this.process) {
      this.process.kill();
    }
  }
}

// Usage example
async function main() {
  const client = new KirhaMCPClient(
    process.env.KIRHA_API_KEY,
    process.env.KIRHA_VERTICAL
  );

  try {
    await client.start();

    // List available tools
    const tools = await client.listTools();
    console.log('Available tools:', tools);

    // Execute a tool
    const result = await client.callTool('example-tool', {
      input: 'Hello, world!'
    });
    console.log('Tool result:', result);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.stop();
  }
}

main();
```

## Python Integration

### Basic Setup

```python
import subprocess
import json
import asyncio
import os
from typing import Dict, Any, List

class KirhaMCPClient:
    def __init__(self, api_key: str, vertical_id: str):
        self.api_key = api_key
        self.vertical_id = vertical_id
        self.process = None
        self.request_id = 1
        self.pending_requests = {}

    async def start(self):
        """Start the Kirha MCP Gateway in stdio mode"""
        env = os.environ.copy()
        env.update({
            'KIRHA_API_KEY': self.api_key,
            'KIRHA_VERTICAL': self.vertical_id,
            'ENABLE_LOGS': 'false'
        })

        self.process = await asyncio.create_subprocess_exec(
            'npx', '@kirha/mcp-server', 'stdio',
            stdin=asyncio.subprocess.PIPE,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            env=env
        )

        # Start reading responses
        asyncio.create_task(self._read_responses())

    async def _read_responses(self):
        """Read and process responses from the gateway"""
        buffer = ""
        
        while True:
            try:
                data = await self.process.stdout.read(1024)
                if not data:
                    break
                    
                buffer += data.decode()
                
                # Process complete JSON messages
                while '\n' in buffer:
                    line, buffer = buffer.split('\n', 1)
                    if line.strip():
                        try:
                            response = json.loads(line)
                            await self._handle_response(response)
                        except json.JSONDecodeError as e:
                            print(f"Failed to parse JSON: {e}, line: {line}")
                            
            except Exception as e:
                print(f"Error reading responses: {e}")
                break

    async def _handle_response(self, response: Dict[str, Any]):
        """Handle a response from the gateway"""
        if 'id' in response and response['id'] in self.pending_requests:
            future = self.pending_requests.pop(response['id'])
            
            if 'error' in response:
                future.set_exception(Exception(response['error']['message']))
            else:
                future.set_result(response.get('result'))

    async def _send_request(self, method: str, params: Dict[str, Any] = None) -> Any:
        """Send a JSON-RPC request and wait for response"""
        request_id = self.request_id
        self.request_id += 1

        request = {
            'jsonrpc': '2.0',
            'id': request_id,
            'method': method
        }
        
        if params:
            request['params'] = params

        # Create future for response
        future = asyncio.Future()
        self.pending_requests[request_id] = future

        # Send request
        request_json = json.dumps(request) + '\n'
        self.process.stdin.write(request_json.encode())
        await self.process.stdin.drain()

        # Wait for response
        return await future

    async def list_tools(self) -> List[Dict[str, Any]]:
        """List all available tools"""
        result = await self._send_request('tools/list')
        return result.get('tools', [])

    async def call_tool(self, name: str, arguments: Dict[str, Any]) -> Dict[str, Any]:
        """Execute a tool with the given arguments"""
        return await self._send_request('tools/call', {
            'name': name,
            'arguments': arguments
        })

    async def stop(self):
        """Stop the gateway process"""
        if self.process:
            self.process.terminate()
            await self.process.wait()

# Usage example
async def main():
    client = KirhaMCPClient(
        api_key=os.getenv('KIRHA_API_KEY'),
        vertical_id=os.getenv('KIRHA_VERTICAL')
    )

    try:
        await client.start()

        # List available tools
        tools = await client.list_tools()
        print(f"Available tools: {len(tools)}")
        for tool in tools:
            print(f"  - {tool['name']}: {tool['description']}")

        # Execute a tool (if available)
        if tools:
            tool_name = tools[0]['name']
            result = await client.call_tool(tool_name, {
                'input': 'Hello, world!'
            })
            print(f"Tool result: {result}")

    except Exception as e:
        print(f"Error: {e}")
    finally:
        await client.stop()

if __name__ == "__main__":
    asyncio.run(main())
```

## HTTP Client Example

For applications that prefer HTTP communication:

### JavaScript (HTTP)

```javascript
import axios from 'axios';

class KirhaHTTPClient {
  constructor(apiKey, verticalId, port = 8022) {
    this.apiKey = apiKey;
    this.verticalId = verticalId;
    this.baseURL = `http://localhost:${port}`;
    this.requestId = 1;
    this.process = null;
  }

  async start() {
    // Start the HTTP server
    const { spawn } = await import('child_process');
    
    this.process = spawn('npx', ['@kirha/mcp-server', 'http'], {
      env: {
        ...process.env,
        KIRHA_API_KEY: this.apiKey,
        KIRHA_VERTICAL: this.verticalId,
        MCP_PORT: '8022'
      }
    });

    // Wait for server to start
    await this.waitForServer();
  }

  async waitForServer(maxRetries = 30) {
    for (let i = 0; i < maxRetries; i++) {
      try {
        await axios.get(`${this.baseURL}/health`);
        return;
      } catch (error) {
        await new Promise(resolve => setTimeout(resolve, 1000));
      }
    }
    throw new Error('Server failed to start');
  }

  async sendRequest(method, params = {}) {
    const request = {
      jsonrpc: '2.0',
      id: this.requestId++,
      method,
      params
    };

    try {
      const response = await axios.post(this.baseURL, request, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.data.error) {
        throw new Error(response.data.error.message);
      }

      return response.data.result;
    } catch (error) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error.message);
      }
      throw error;
    }
  }

  async listTools() {
    return this.sendRequest('tools/list');
  }

  async callTool(name, arguments) {
    return this.sendRequest('tools/call', { name, arguments });
  }

  async stop() {
    if (this.process) {
      this.process.kill();
    }
  }
}

// Usage
const client = new KirhaHTTPClient(
  process.env.KIRHA_API_KEY,
  process.env.KIRHA_VERTICAL
);

await client.start();
const tools = await client.listTools();
console.log('Tools:', tools);
await client.stop();
```

## Go Integration

For Go applications that need to embed the gateway:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "bufio"
    "strings"
    "sync"
)

type MCPClient struct {
    cmd       *exec.Cmd
    stdin     io.WriteCloser
    stdout    io.ReadCloser
    requestID int
    pending   map[int]chan *MCPResponse
    mu        sync.Mutex
}

type MCPRequest struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      int         `json:"id"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      int         `json:"id"`
    Result  interface{} `json:"result,omitempty"`
    Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func NewMCPClient(apiKey, verticalID string) *MCPClient {
    return &MCPClient{
        pending: make(map[int]chan *MCPResponse),
    }
}

func (c *MCPClient) Start(ctx context.Context) error {
    c.cmd = exec.CommandContext(ctx, "npx", "@kirha/mcp-server", "stdio")
    c.cmd.Env = append(os.Environ(),
        fmt.Sprintf("KIRHA_API_KEY=%s", apiKey),
        fmt.Sprintf("KIRHA_VERTICAL=%s", verticalID),
        "ENABLE_LOGS=false",
    )

    stdin, err := c.cmd.StdinPipe()
    if err != nil {
        return err
    }
    c.stdin = stdin

    stdout, err := c.cmd.StdoutPipe()
    if err != nil {
        return err
    }
    c.stdout = stdout

    if err := c.cmd.Start(); err != nil {
        return err
    }

    // Start reading responses
    go c.readResponses()

    return nil
}

func (c *MCPClient) readResponses() {
    scanner := bufio.NewScanner(c.stdout)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        var response MCPResponse
        if err := json.Unmarshal([]byte(line), &response); err != nil {
            fmt.Printf("Failed to parse response: %v\n", err)
            continue
        }

        c.mu.Lock()
        if ch, exists := c.pending[response.ID]; exists {
            delete(c.pending, response.ID)
            c.mu.Unlock()
            ch <- &response
        } else {
            c.mu.Unlock()
        }
    }
}

func (c *MCPClient) sendRequest(method string, params interface{}) (*MCPResponse, error) {
    c.mu.Lock()
    id := c.requestID
    c.requestID++
    
    responseCh := make(chan *MCPResponse, 1)
    c.pending[id] = responseCh
    c.mu.Unlock()

    request := MCPRequest{
        JSONRPC: "2.0",
        ID:      id,
        Method:  method,
        Params:  params,
    }

    data, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }

    _, err = c.stdin.Write(append(data, '\n'))
    if err != nil {
        return nil, err
    }

    response := <-responseCh
    if response.Error != nil {
        return nil, fmt.Errorf("MCP error: %s", response.Error.Message)
    }

    return response, nil
}

func (c *MCPClient) ListTools() (interface{}, error) {
    response, err := c.sendRequest("tools/list", nil)
    if err != nil {
        return nil, err
    }
    return response.Result, nil
}

func (c *MCPClient) CallTool(name string, arguments map[string]interface{}) (interface{}, error) {
    params := map[string]interface{}{
        "name":      name,
        "arguments": arguments,
    }
    
    response, err := c.sendRequest("tools/call", params)
    if err != nil {
        return nil, err
    }
    return response.Result, nil
}

func (c *MCPClient) Stop() error {
    if c.stdin != nil {
        c.stdin.Close()
    }
    if c.cmd != nil {
        return c.cmd.Wait()
    }
    return nil
}

// Usage example
func main() {
    client := NewMCPClient(
        os.Getenv("KIRHA_API_KEY"),
        os.Getenv("KIRHA_VERTICAL"),
    )

    ctx := context.Background()
    if err := client.Start(ctx); err != nil {
        log.Fatal(err)
    }
    defer client.Stop()

    // List tools
    tools, err := client.ListTools()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tools: %+v\n", tools)

    // Call a tool
    result, err := client.CallTool("example-tool", map[string]interface{}{
        "input": "Hello, world!",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Result: %+v\n", result)
}
```

## Best Practices

### Error Handling

```javascript
// Always handle errors gracefully
try {
  const result = await client.callTool('tool-name', args);
  console.log('Success:', result);
} catch (error) {
  if (error.message.includes('TOOL_NOT_FOUND')) {
    console.log('Tool not available, using fallback...');
  } else if (error.message.includes('TIMEOUT')) {
    console.log('Operation timed out, retrying...');
  } else {
    console.error('Unexpected error:', error);
  }
}
```

### Resource Management

```python
# Use context managers for proper cleanup
from contextlib import asynccontextmanager

@asynccontextmanager
async def kirha_client(api_key, vertical_id):
    client = KirhaMCPClient(api_key, vertical_id)
    try:
        await client.start()
        yield client
    finally:
        await client.stop()

# Usage
async with kirha_client(api_key, vertical_id) as client:
    tools = await client.list_tools()
    # Client is automatically stopped when exiting the context
```

### Performance Optimization

```javascript
// Cache tool schemas to avoid repeated list_tools calls
class CachedKirhaClient extends KirhaMCPClient {
  constructor(apiKey, verticalId) {
    super(apiKey, verticalId);
    this.toolsCache = null;
    this.cacheExpiry = null;
  }

  async getTools() {
    const now = Date.now();
    if (!this.toolsCache || now > this.cacheExpiry) {
      this.toolsCache = await this.listTools();
      this.cacheExpiry = now + (5 * 60 * 1000); // 5 minutes
    }
    return this.toolsCache;
  }
}
```

## Troubleshooting

### Common Issues

1. **Process spawn errors**: Ensure `@kirha/mcp-server` is installed globally
2. **JSON parsing errors**: Check for mixed stdout/stderr output
3. **Timeout issues**: Increase timeout values for slow tools
4. **Authentication errors**: Verify API key and vertical ID

### Debug Mode

Enable debugging by setting environment variables:

```javascript
const client = new KirhaMCPClient(apiKey, verticalId);
// Enable logs for debugging
process.env.ENABLE_LOGS = 'true';
await client.start();
```

For more examples and support, see the [API Documentation](../API.md).