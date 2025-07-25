#!/usr/bin/env python3
"""
MCP Server using FastMCP
Provides tools for basic math operations and text processing
"""

from fastapi import FastAPI, status
from fastmcp import FastMCP
from pydantic import BaseModel, Field
import json
import uvicorn
import threading


# ------------------------------------------------ FastAPI Setup ----------------------------------------------
app = FastAPI()

class HealthCheck(BaseModel):
    """Response model to validate and return when performing a health check."""

    status: str = "OK"


@app.get(
    "/health",
    tags=["healthcheck"],
    summary="Perform a Health Check",
    response_description="Return HTTP Status Code 200 (OK)",
    status_code=status.HTTP_200_OK,
    response_model=HealthCheck,
)
def get_health() -> HealthCheck:
    """
    ## Perform a Health Check
    Endpoint to perform a healthcheck on. This endpoint can primarily be used Docker
    to ensure a robust container orchestration and management is in place. Other
    services which rely on proper functioning of the API service will not deploy if this
    endpoint returns any other HTTP status code except 200 (OK).
    Returns:
        HealthCheck: Returns a JSON response with the health status
    """
    return HealthCheck(status="OK")

# ---------------------------------------------- MCP Tools ----------------------------------------------

# Initialize FastMCP server
mcp = FastMCP("Math and Text Server")
@mcp.tool()

async def add_numbers(a: float, b: float) -> str:
    """Add two numbers together"""
    result = a + b
    return f"The sum of {a} and {b} is {result}"


@mcp.tool()
async def multiply_numbers(a: float, b: float) -> str:
    """Multiply two numbers together"""
    result = a * b
    return f"The product of {a} and {b} is {result}"


@mcp.tool()
async def process_text(text: str, operation: str) -> str:
    """Process text with various operations"""
    text = text
    operation = operation.lower()

    if operation == "upper":
        result = text.upper()
        return f"Uppercase: {result}"
    elif operation == "lower":
        result = text.lower()
        return f"Lowercase: {result}"
    elif operation == "reverse":
        result = text[::-1]
        return f"Reversed: {result}"
    else:
        return f"Unknown operation: {operation}. Available: upper, lower, reverse"


def main():
    """Run the MCP server with its JSON-RPC HTTP transport"""
    print("Starting MCP server...")
    print("Available tools:")
    print("  - add_numbers: Add two numbers together")
    print("  - multiply_numbers: Multiply two numbers together")
    print("  - process_text: Process text with various operations")

    # Use the FastMCP's run method which handles the transport automatically
    # This is the correct way to start a FastMCP server
    #mcp.run(transport="http", host="0.0.0.0", port=9000)

     # Start FastAPI server in a separate thread
    def start_fastapi():
        uvicorn.run(app, host="0.0.0.0", port=9001)
    
    fastapi_thread = threading.Thread(target=start_fastapi, daemon=True)
    fastapi_thread.start()
    
    print("FastAPI health endpoint started on http://0.0.0.0:9001/health")
    
    # Start the MCP server (this blocks)
    print("Starting MCP JSON-RPC server on http://0.0.0.0:9000")
    mcp.run(transport="http", host="0.0.0.0", port=9000)

if __name__ == "__main__":
    main()