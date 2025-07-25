package mcpserver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

// MCPServer wraps the underlying MCP server implementation.
type MCPServer struct {
	mcpSrv *server.MCPServer

	storage map[string]User // In-memory storage for users
}

// NewMCPServer initializes a new MCP server with RESTful tools.
// It sets up tools for getting user information, creating a user, and searching users.
// Each tool has its own handler function that processes the request and returns a result.
func NewMCPServer() (MCPServer, error) {
	s := server.NewMCPServer("StreamableHTTP API Server", "1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	mcpServer := MCPServer{mcpSrv: s}
	mcpServer.populate() // Populate the in-memory storage with initial users

	// Add RESTful tools
	s.AddTool(
		mcp.NewTool("get_user",
			mcp.WithDescription("Get user information"),
			mcp.WithString("user_id", mcp.Required()),
		),
		buildGetUserHandler(&mcpServer),
	)

	s.AddTool(
		mcp.NewTool("create_user",
			mcp.WithDescription("Create a new user"),
			mcp.WithString("user_id", mcp.Required()),
			mcp.WithString("name", mcp.Required()),
			mcp.WithString("email", mcp.Required()),
			mcp.WithNumber("age", mcp.Min(0)),
		),
		buildCreateUserHandler(&mcpServer),
	)

	return mcpServer, nil
}

// Start initializes the MCP server and starts the StreamableHTTP server on port 8080.
// It listens for incoming requests and handles tool calls as defined in the MCP server.
func (s MCPServer) Start(address string) error {
	// start the MCP server spec handler
	go startOpenAPIServer()

	// Start StreamableHTTP server
	log.Println("Starting StreamableHTTP server on :8090")
	httpServer := server.NewStreamableHTTPServer(s.mcpSrv,
		server.WithHeartbeatInterval(30*time.Second),
		server.WithStateLess(true),
	)

	if err := httpServer.Start(address); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *MCPServer) readUser(userID string) (*User, error) {
	// Simulate a database lookup
	user, exists := s.storage[userID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	return &user, nil
}

func buildGetUserHandler(s *MCPServer) server.ToolHandlerFunc {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID := req.GetString("user_id", "")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}

		// Simulate database lookup
		user, err := s.readUser(userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %s", userID)
		}

		// Return the user information as a text result
		// This is a simple example; you might want to return structured data in a real application
		// For example, you could return a JSON object with user details
		// Here we just return a simple text response for demonstration purposes.
		return mcp.NewToolResultText(fmt.Sprintf("User found: %s, Email: %s, Age: %d", user.Name, user.Email, user.Age)), nil
	}
}

// createUser creates a new user with the provided name, email, and age.
// It validates the input, generates a unique ID, and saves the user to the in-memory storage.
// If successful, it returns the created user; otherwise, it returns an error.
func (s *MCPServer) createUser(_ context.Context, userID, name, email string, age int) error {
	if _, ok := s.storage[userID]; ok {
		return fmt.Errorf("user with ID %s already exists", userID)
	}

	// Validate input
	if name == "" || email == "" {
		return fmt.Errorf("name and email are required")
	}
	if age < 0 {
		return fmt.Errorf("age must be a non-negative integer")
	}

	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	// Create user
	user := User{
		ID:        userID,
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: time.Now(),
	}

	// Save user to in-memory storage
	s.storage[user.ID] = user

	// Log the creation
	log.Printf("User created: ID=%s, Name=%s, Email=%s, Age=%d", user.ID, user.Name, user.Email, user.Age)

	return nil
}

func buildCreateUserHandler(s *MCPServer) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID := req.GetString("user_id", "")
		name := req.GetString("name", "")
		email := req.GetString("email", "")
		age := req.GetInt("age", 0)

		if age < 0 {
			return nil, fmt.Errorf("age must be a non-negative integer")
		}

		if !isValidEmail(email) {
			return nil, fmt.Errorf("invalid email format: %s", email)
		}

		if err := s.createUser(ctx, userID, name, email, age); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		return mcp.NewToolResultText("OK"), nil
	}
}

func isValidEmail(email string) bool {
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (s *MCPServer) populate() {
	// Populate the in-memory storage with some initial users
	s.storage = make(map[string]User)
	s.storage["1"] = User{ID: "1", Name: "Alice", Email: "alice@email.com", Age: 28, CreatedAt: time.Now()}
	s.storage["2"] = User{ID: "2", Name: "Bob", Email: "bob@email.com", Age: 32, CreatedAt: time.Now()}
	s.storage["3"] = User{ID: "3", Name: "Charlie", Email: "charlie@email.com", Age: 22, CreatedAt: time.Now()}
	log.Println("In-memory storage populated with initial users")
}
