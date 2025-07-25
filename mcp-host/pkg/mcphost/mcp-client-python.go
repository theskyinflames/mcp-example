package mcphost

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewPythonClient creates a new MCP client that connects to the Python MCP server.
// This client will communicate with the Python server using the JSON-RPC protocol.
func NewPythonClient(ctx context.Context, URL string) (*client.Client, error) {
	fmt.Println("Initializing HTTP client for Python server...")

	// The Python server's JSON-RPC transport listens on the /mcp path by default
	httpTransport, err := transport.NewStreamableHTTP(URL + "/mcp")
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP transport for Python server: %w", err)
	}

	// Create client with the transport
	c := client.NewClient(httpTransport)

	// Start the client
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start Python client: %w", err)
	}

	// Set up notification handler
	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification from Python server: %s\n", notification.Method)
	})

	// Initialize the client
	fmt.Println("Initializing client for Python server...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Python Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize with Python server: %v", err)
	}

	// Display server information
	fmt.Printf("Connected to Python server: %s (version %s)\n",
		serverInfo.ServerInfo.Name,
		serverInfo.ServerInfo.Version)
	fmt.Printf("Python server capabilities: %+v\n", serverInfo.Capabilities)

	return c, nil
}

// MCPToolsSchemaPythonSrv retrieves the list of tools from the Python MCP server.
func MCPToolsSchemaPythonSrv(ctx context.Context, mcpClient *client.Client) (string, error) {
	// Get the list of resources from the Python MCP server
	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return "", err
	}

	var response []*mcp.Tool
	for _, t := range tools.Tools {
		response = append(response, &t)
	}

	pythonMCPSrvTools, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Python tools: %v", err)
	}

	return string(pythonMCPSrvTools), nil
}
