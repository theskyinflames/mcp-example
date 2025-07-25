package main

import (
	"log"
	"os"

	"mcp-server-go/internal/mcpserver"
)

const (
	// MCPServerAddress is the address of the MCP server
	MCPServerAddress = ":8090"
)

func main() {
	// Initialize the MCP server with RESTful tools
	s, err := mcpserver.NewMCPServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	if err := s.Start(mcpServerAddress()); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

func mcpServerAddress() string {
	addr := os.Getenv("MCP_SERVER_ENDPOINT")
	if addr == "" {
		return MCPServerAddress
	}
	return addr
}
