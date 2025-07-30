package mcphost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type (
	// ChatRequest represents the request structure for the LLM chat API.
	ChatRequest struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
		Stream   bool          `json:"stream,omitempty"` // Optional, set to true if you want streaming responses
	}

	// ChatMessage represents a single message in the chat conversation.
	ChatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// ChatResponse represents the response structure from the LLM chat API.
	ChatResponse struct {
		Choices []struct {
			Message ChatMessage `json:"message"`
		} `json:"choices"`
	}
)

// LLMClient is a client for interacting with the Deepseek LLM API.
type LLMClient struct {
	apiKey     string
	llmURL     string
	systemRole string // This can be used to set the system role for the LLM
}

// NewLLMClient creates a new LLM client with the provided Deepseek API key.
func NewLLMClient(apiKey, llmURL string, goMCPSrvTools string, pythonMCPSrvTools string) (LLMClient, error) {
	fmt.Println("Available tools from Go MCP server:", goMCPSrvTools)
	fmt.Println("Available tools from Python MCP server:", pythonMCPSrvTools)

	systemRole := fmt.Sprintf(`You are an assistant that can use tools to answer questions.
When a tool is needed, you must respond with a JSON object containing the tool name and its inputs.
Use the tag "tool" to indicate the tool name and "inputs" for the arguments.
Do not add any other text to the response.

Available tools from Go MCP server: %s\n,
Available tools from Python MCP server: %s\n 

Make sure the tool's request is valid and 
that it follows the above description for each MCP server`, goMCPSrvTools, pythonMCPSrvTools)

	return LLMClient{
		apiKey:     apiKey,
		llmURL:     llmURL,
		systemRole: systemRole, // Store the system role for use in requests
	}, nil
}

// CallLLM sends a user inpu to the LLM and returns the response as a string.
func (c LLMClient) CallLLM(userInput string) (string, error) {
	reqBody := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "system", Content: c.systemRole},
			{Role: "user", Content: userInput},
		},
		Stream: false, // Set to true if you want streaming responses
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.llmURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM call failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed ChatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	planJSON := parsed.Choices[0].Message.Content

	// Normalize the plan JSON if needed
	normalizedPlan, err := normalizePlanFrmLLM(planJSON)
	if err != nil {
		return "", fmt.Errorf("failed to normalize plan JSON: %v", err)
	}
	return normalizedPlan, nil
}

func normalizePlanFrmLLM(planJSON string) (string, error) {
	// Clean the JSON string if it's wrapped in Markdown code block
	if strings.HasPrefix(planJSON, "```json\n") {
		planJSON = strings.TrimPrefix(planJSON, "```json\n")
		planJSON = strings.TrimSuffix(planJSON, "\n```")
	} else if strings.HasPrefix(planJSON, "```") { // Handle cases with just ```
		planJSON = strings.TrimPrefix(planJSON, "```")
		planJSON = strings.TrimSuffix(planJSON, "```")
	}
	return strings.TrimSpace(planJSON), nil
}
