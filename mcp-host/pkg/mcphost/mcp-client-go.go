package mcphost

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP client example https://github.com/mark3labs/mcp-go/blob/main/examples/simple_client/main.go#L1

// NewClient creates a new MCP client that connects to the MCP server.
// Available options: https://github.com/mark3labs/mcp-go/blob/main/client/transport/streamable_http.go#L22
// http.Client, headers, auth, ...
func NewClient(ctx context.Context, URL string) (*client.Client, error) {
	fmt.Println("Initializing HTTP client...")
	// Create HTTP transport
	httpTransport, err := transport.NewStreamableHTTP(URL + "/mcp")
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP transport: %v", err)
	}

	// Create client with the transport
	c := client.NewClient(httpTransport)

	// Start the client
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("Failed to start client: %v", err)
	}

	// Set up notification handler
	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Display server information
	fmt.Printf("Connected to server: %s (version %s)\n",
		serverInfo.ServerInfo.Name,
		serverInfo.ServerInfo.Version)
	fmt.Printf("Server capabilities: %+v\n", serverInfo.Capabilities)

	return c, nil
}

// MCPToolsSchemaGoSrv retrieves the list of available tools from the MCP server.
func MCPToolsSchemaGoSrv(ctx context.Context, c *client.Client) (string, error) {
	tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return "", err
	}

	var response []*mcp.Tool
	for _, t := range tools.Tools {
		response = append(response, &t)
	}

	goTools, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tools: %v", err)
	}

	return string(goTools), nil
}
