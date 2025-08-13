#!/bin/bash

set -e

# get the API key from the environment variable
API_KEY=${DEEPSEEK_API_KEY}

# run the Go binary with the API key
if [ -z "$API_KEY" ]; then
  echo "Error: DEEPSEEK_API_KEY is not set."
  exit 1
fi 

echo "Starting MCP servers in background..."
# Start the servers first
DEEPSEEK_API_KEY="$API_KEY" docker-compose up -d mcp-server-go mcp-server-python

echo "Waiting for servers to be ready..."
sleep 5

echo "Starting interactive MCP host..."
# Run the host interactively
DEEPSEEK_API_KEY="$API_KEY" docker-compose run --rm mcp-host

# Clean up
echo "Stopping servers..."
docker-compose down
