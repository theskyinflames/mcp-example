package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"mcp-example-host/pkg/mcphost"

	"github.com/mark3labs/mcp-go/client"
)

const (

	// GoMCPServerAddress s the address of the MCP server
	GoMCPServerAddress = "http://0.0.0.0:8090"

	// PythonMCPServerAddress is the address of the Python MCP server
	PythonMCPServerAddress = "http://0.0.0.0:9000"

	// LLMServerAddress is the address of the LLM server
	LLMServerAddress = "https://api.deepseek.com/v1/chat/completions"
)

func main() {
	ctx := context.Background()

	// Create a new MCP client for the Go server
	goMCPClient, err := buildMCPClient(ctx, goMCPServerAddr())
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}

	// Create a new MCP client for the Python server
	pythonMCPClient, err := buildMCPClient(ctx, pythonMCPServerAddr())
	if err != nil {
		log.Fatalf("Failed to create Python MCP client: %v", err)
	}

	// Create a new LLM client
	// This client will connect to the Deepseek LLM API
	// and use the tools provided by the MCP server
	llmClient, err := buildLLMClient(ctx, goMCPClient, pythonMCPClient)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	// Create a new LLM application
	// This application integrates the MCP client and LLM client
	// It will use the LLM to analyze user queries and call the appropriate tools
	// The LLM application will also handle the tool calling logic
	// and manage the interaction between the user, LLM, and tools.
	llmApp, err := mcphost.NewMCPHost(goMCPClient, pythonMCPClient, llmClient)
	if err != nil {
		log.Fatalf("Failed to create LLM application: %v", err)
	}

	if err := runExamples(ctx, llmApp); err != nil {
		log.Fatalf("Failed to run examples: %v", err)
	}
}

func runExamples(ctx context.Context, llmApp *mcphost.MCPHost) error {
	// Run a tool example to read user data
	q := "Read user with Id 2"
	fmt.Println("\n\nRunning user query:", q)
	response, err := llmApp.RunUserQuery(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to run user query: %w", err)
	}

	// Print the response from the tool
	fmt.Println("Read user response:", response)

	// Run another tool example to create a new user
	userID := time.Now().UnixMilli()
	q = fmt.Sprintf("Create a new user with Id %d, name 'John Doe', email 'jhondoe@email.com', age 30", userID)
	fmt.Println("\n\nRunning user query:", q)
	response, err = llmApp.RunUserQuery(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to run user query: %w", err)
	}

	// Print the response from the tool
	fmt.Println("Create user response:", response)

	// Run another tool example to add two numbers
	q = "Add two numbers 5 and 10"
	fmt.Println("\n\nRunning user query:", q)
	response, err = llmApp.RunUserQuery(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to run user query: %w", err)
	}
	// Print the response from the tool
	fmt.Println("Add numbers response:", response)

	return nil
}

func buildMCPClient(ctx context.Context, URL string) (*client.Client, error) {
	var (
		MCPClient *client.Client
		attempts  int
		err       error
	)

	const maxAttempts = 5

	for {
		fmt.Printf("Retrying to create Python MCP client (%d/%d)...", attempts, maxAttempts)
		MCPClient, err = mcphost.NewPythonClient(ctx, URL)
		if err == nil {
			break
		}
		if attempts == maxAttempts {
			break // Exit after 3 attempts
		}
		attempts++
		time.Sleep(time.Second) // Wait before retrying
	}

	return MCPClient, err
}

// buildLLMClient initializes a new LLM client with the MCP client
// It retrieves the available tools from the MCP server and uses them to configure the LLM client.
func buildLLMClient(ctx context.Context, goMCPClient *client.Client, pythonMCPClient *client.Client) (mcphost.LLMClient, error) {
	// Get client toolsGo
	toolsGo, err := mcphost.MCPToolsSchemaGoSrv(ctx, goMCPClient)
	if err != nil {
		return mcphost.LLMClient{}, fmt.Errorf("failed to list tools: %w", err)
	}

	// Get python client tools
	toolsPython, err := mcphost.MCPToolsSchemaPythonSrv(ctx, pythonMCPClient)
	if err != nil {
		return mcphost.LLMClient{}, fmt.Errorf("failed to list tools from Python server: %w", err)
	}

	apiKey, err := deepseekAPIKey()
	if err != nil {
		return mcphost.LLMClient{}, fmt.Errorf("failed to get Deepseek API key: %w", err)
	}

	return mcphost.NewLLMClient(
		apiKey,
		LLMServerAddress,
		toolsGo,
		toolsPython,
	)
}

// deepseekAPIKey get the API key for Deepseek from environment
func deepseekAPIKey() (string, error) {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY environment variable is not set")
	}

	fmt.Printf("Using Deepseek API key: %s\n", key)

	return key, nil
}

func goMCPServerAddr() string {
	addr := os.Getenv("USERS_MCP_SERVER_ADDR")
	if addr == "" {
		return GoMCPServerAddress
	}

	fmt.Printf("Using MCP server address: %s\n", addr)

	return addr
}

func pythonMCPServerAddr() string {
	addr := os.Getenv("CALC_MCP_SERVER_ADDR")
	if addr == "" {
		return PythonMCPServerAddress
	}

	fmt.Printf("Using Python MCP server address: %s\n", addr)

	return addr
}
