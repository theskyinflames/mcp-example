package mcphost

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	// MCPHost represents the main application that integrates the MCP client and LLM client.
	MCPHost struct {
		mcpClient       *client.Client
		mcpClientPython *client.Client
		llmClient       LLMClient
	}
)

// LLMToolCallPlan defines the structure for the LLM's plan to call a tool.
// We are assuming the LLM will respond with a JSON in this format.
type LLMToolCallPlan struct {
	ToolName  string                 `json:"tool"`
	Arguments map[string]interface{} `json:"inputs"`
}

// NewMCPHost initializes a new LLM application with the MCP client and LLM client.
// It connects to the MCP server at the specified address and retrieves available tools.
// It also initializes the LLM client with the Deepseek API key and system role that includes the available tools.
// The system role is set to inform the LLM about the available tools and how to use them.
// The LLM client is configured to use the Deepseek API for chat completions.
func NewMCPHost(mcpClient *client.Client, mcpClientPython *client.Client, llmClient LLMClient) (*MCPHost, error) {
	return &MCPHost{
		mcpClient:       mcpClient,
		mcpClientPython: mcpClientPython,
		llmClient:       llmClient,
	}, nil
}

// RunUserQuery demonstrates how to use the LLM to call a tool based on user input.
// It simulates a user query, lets the LLM analyze it, and then calls the appropriate tool.
// This is a simple example that assumes the LLM will return a JSON object with the tool name and arguments.
// It also includes error handling for cases where the LLM does not return a valid tool call plan.
func (l *MCPHost) RunUserQuery(ctx context.Context, query string) (string, error) {
	// Step 1: LLM analysis
	planJSON, err := l.llmClient.CallLLM(query)
	if err != nil {
		log.Fatalf("LLM call failed: %v", err)
	}
	log.Println("LLM output (raw JSON):", planJSON)

	// Step 2: Parse the LLM's plan
	var plan LLMToolCallPlan
	if err := json.Unmarshal([]byte(planJSON), &plan); err != nil {
		// This is a fallback if the LLM does not return a valid JSON for tool calling
		// or if the response is a direct answer instead of a tool call plan.
		return "", fmt.Errorf("Failed to parse LLM plan from JSON: %v", err) // Exit if we can't parse a tool call
	}

	// Check if a tool name was actually provided by the LLM
	if plan.ToolName == "" {
		return "", fmt.Errorf("LLM did not specify a tool to call. LLM response: %v", planJSON)
	}

	log.Printf("LLM identified tool: %s with arguments: %v", plan.ToolName, plan.Arguments)

	// Step 3: Call tool with MCP based on LLM's plan
	request := mcp.CallToolRequest{}
	request.Params.Name = plan.ToolName
	request.Params.Arguments = plan.Arguments

	// Try the Go MCP server first
	response, err := l.mcpClient.CallTool(ctx, request)
	if err != nil {
		// If the tool is not found in the Go server, try the Python server
		log.Printf("Tool '%s' not found in Go server, trying Python server: %v", plan.ToolName, err)
		response, err = l.mcpClientPython.CallTool(ctx, request)
		if err != nil {
			return "", fmt.Errorf("Tool call failed on both servers: %v", err)
		}
	}

	if response.IsError {
		return "", fmt.Errorf("Tool call error: %s", response.Content)
	}

	textResponse, ok := mcp.AsTextContent(response.Content[0])
	if !ok {
		return "", fmt.Errorf("Tool call did not return text content: %v", response.Content)
	}

	return textResponse.Text, nil
}
