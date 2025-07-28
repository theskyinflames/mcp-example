#!/bin/bash

set -e

# get the API key from the environment variable
API_KEY=${DEEPSEEK_API_KEY}

# run the Go binary with the API key
if [ -z "$API_KEY" ]; then
  echo "Error: DEEPSEEK_API_KEY is not set."
  exit 1
fi 

# Run the MCP demo
DEEPSEEK_API_KEY="$API_KEY" docker-compose up --build --abort-on-container-exit
