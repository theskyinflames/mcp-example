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


@app.get(
    "/mcp/openapi.json",
    tags=["mcp"],
    summary="Get MCP Tools OpenAPI Specification",
    response_description="Returns OpenAPI 3.0 specification for MCP tools",
)
def get_mcp_openapi_spec():
    """
    ## Get MCP Tools OpenAPI Specification
    Returns an OpenAPI 3.0 specification that documents all available MCP tools
    as REST endpoints. This provides a standard way to understand the available
    tools and their parameters.
    """
    openapi_spec = {
        "openapi": "3.0.0",
        "info": {
            "title": "MCP Math and Text Server Tools",
            "description": "Model Context Protocol server providing mathematical and text processing tools",
            "version": "1.12.0",
            "contact": {
                "name": "MCP Example Project",
                "url": "https://github.com/theskyinflames/mcp-example",
            },
        },
        "servers": [
            {"url": "http://localhost:9000", "description": "MCP JSON-RPC Server"}
        ],
        "paths": {
            "/message": {
                "post": {
                    "tags": ["MCP Tools"],
                    "summary": "Call MCP Tools",
                    "description": "JSON-RPC endpoint for calling MCP tools",
                    "requestBody": {
                        "required": True,
                        "content": {
                            "application/json": {
                                "schema": {
                                    "oneOf": [
                                        {
                                            "type": "object",
                                            "title": "Add Numbers Tool",
                                            "properties": {
                                                "jsonrpc": {
                                                    "type": "string",
                                                    "enum": ["2.0"],
                                                },
                                                "id": {"type": "integer"},
                                                "method": {
                                                    "type": "string",
                                                    "enum": ["tools/call"],
                                                },
                                                "params": {
                                                    "type": "object",
                                                    "properties": {
                                                        "name": {
                                                            "type": "string",
                                                            "enum": ["add_numbers"],
                                                        },
                                                        "arguments": {
                                                            "type": "object",
                                                            "properties": {
                                                                "a": {
                                                                    "type": "number",
                                                                    "description": "First number to add",
                                                                },
                                                                "b": {
                                                                    "type": "number",
                                                                    "description": "Second number to add",
                                                                },
                                                            },
                                                            "required": ["a", "b"],
                                                        },
                                                    },
                                                    "required": ["name", "arguments"],
                                                },
                                            },
                                            "required": [
                                                "jsonrpc",
                                                "id",
                                                "method",
                                                "params",
                                            ],
                                        },
                                        {
                                            "type": "object",
                                            "title": "Multiply Numbers Tool",
                                            "properties": {
                                                "jsonrpc": {
                                                    "type": "string",
                                                    "enum": ["2.0"],
                                                },
                                                "id": {"type": "integer"},
                                                "method": {
                                                    "type": "string",
                                                    "enum": ["tools/call"],
                                                },
                                                "params": {
                                                    "type": "object",
                                                    "properties": {
                                                        "name": {
                                                            "type": "string",
                                                            "enum": [
                                                                "multiply_numbers"
                                                            ],
                                                        },
                                                        "arguments": {
                                                            "type": "object",
                                                            "properties": {
                                                                "a": {
                                                                    "type": "number",
                                                                    "description": "First number to multiply",
                                                                },
                                                                "b": {
                                                                    "type": "number",
                                                                    "description": "Second number to multiply",
                                                                },
                                                            },
                                                            "required": ["a", "b"],
                                                        },
                                                    },
                                                    "required": ["name", "arguments"],
                                                },
                                            },
                                            "required": [
                                                "jsonrpc",
                                                "id",
                                                "method",
                                                "params",
                                            ],
                                        },
                                        {
                                            "type": "object",
                                            "title": "Process Text Tool",
                                            "properties": {
                                                "jsonrpc": {
                                                    "type": "string",
                                                    "enum": ["2.0"],
                                                },
                                                "id": {"type": "integer"},
                                                "method": {
                                                    "type": "string",
                                                    "enum": ["tools/call"],
                                                },
                                                "params": {
                                                    "type": "object",
                                                    "properties": {
                                                        "name": {
                                                            "type": "string",
                                                            "enum": ["process_text"],
                                                        },
                                                        "arguments": {
                                                            "type": "object",
                                                            "properties": {
                                                                "text": {
                                                                    "type": "string",
                                                                    "description": "Text to process",
                                                                },
                                                                "operation": {
                                                                    "type": "string",
                                                                    "enum": [
                                                                        "upper",
                                                                        "lower",
                                                                        "reverse",
                                                                    ],
                                                                    "description": "Operation to perform on the text",
                                                                },
                                                            },
                                                            "required": [
                                                                "text",
                                                                "operation",
                                                            ],
                                                        },
                                                    },
                                                    "required": ["name", "arguments"],
                                                },
                                            },
                                            "required": [
                                                "jsonrpc",
                                                "id",
                                                "method",
                                                "params",
                                            ],
                                        },
                                        {
                                            "type": "object",
                                            "title": "List Tools",
                                            "properties": {
                                                "jsonrpc": {
                                                    "type": "string",
                                                    "enum": ["2.0"],
                                                },
                                                "id": {"type": "integer"},
                                                "method": {
                                                    "type": "string",
                                                    "enum": ["tools/list"],
                                                },
                                                "params": {"type": "object"},
                                            },
                                            "required": [
                                                "jsonrpc",
                                                "id",
                                                "method",
                                                "params",
                                            ],
                                        },
                                    ]
                                },
                                "examples": {
                                    "add_numbers_example": {
                                        "summary": "Add two numbers",
                                        "value": {
                                            "jsonrpc": "2.0",
                                            "id": 1,
                                            "method": "tools/call",
                                            "params": {
                                                "name": "add_numbers",
                                                "arguments": {"a": 5, "b": 10},
                                            },
                                        },
                                    },
                                    "multiply_numbers_example": {
                                        "summary": "Multiply two numbers",
                                        "value": {
                                            "jsonrpc": "2.0",
                                            "id": 2,
                                            "method": "tools/call",
                                            "params": {
                                                "name": "multiply_numbers",
                                                "arguments": {"a": 6, "b": 7},
                                            },
                                        },
                                    },
                                    "process_text_example": {
                                        "summary": "Process text (uppercase)",
                                        "value": {
                                            "jsonrpc": "2.0",
                                            "id": 3,
                                            "method": "tools/call",
                                            "params": {
                                                "name": "process_text",
                                                "arguments": {
                                                    "text": "hello world",
                                                    "operation": "upper",
                                                },
                                            },
                                        },
                                    },
                                    "list_tools_example": {
                                        "summary": "List available tools",
                                        "value": {
                                            "jsonrpc": "2.0",
                                            "id": 4,
                                            "method": "tools/list",
                                            "params": {},
                                        },
                                    },
                                },
                            }
                        },
                    },
                    "responses": {
                        "200": {
                            "description": "Successful tool execution",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "jsonrpc": {
                                                "type": "string",
                                                "enum": ["2.0"],
                                            },
                                            "id": {"type": "integer"},
                                            "result": {
                                                "oneOf": [
                                                    {
                                                        "type": "object",
                                                        "title": "Tool Call Result",
                                                        "properties": {
                                                            "content": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "properties": {
                                                                        "type": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "text"
                                                                            ],
                                                                        },
                                                                        "text": {
                                                                            "type": "string"
                                                                        },
                                                                    },
                                                                },
                                                            },
                                                            "isError": {
                                                                "type": "boolean"
                                                            },
                                                        },
                                                    },
                                                    {
                                                        "type": "object",
                                                        "title": "Tools List Result",
                                                        "properties": {
                                                            "tools": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string"
                                                                        },
                                                                        "description": {
                                                                            "type": "string"
                                                                        },
                                                                        "inputSchema": {
                                                                            "type": "object"
                                                                        },
                                                                    },
                                                                },
                                                            }
                                                        },
                                                    },
                                                ]
                                            },
                                        },
                                    },
                                    "examples": {
                                        "add_result": {
                                            "summary": "Addition result",
                                            "value": {
                                                "jsonrpc": "2.0",
                                                "id": 1,
                                                "result": {
                                                    "content": [
                                                        {
                                                            "type": "text",
                                                            "text": "The sum of 5 and 10 is 15",
                                                        }
                                                    ],
                                                    "isError": False,
                                                },
                                            },
                                        },
                                        "tools_list": {
                                            "summary": "Available tools",
                                            "value": {
                                                "jsonrpc": "2.0",
                                                "id": 4,
                                                "result": {
                                                    "tools": [
                                                        {
                                                            "name": "add_numbers",
                                                            "description": "Add two numbers together",
                                                            "inputSchema": {
                                                                "type": "object",
                                                                "properties": {
                                                                    "a": {
                                                                        "type": "number"
                                                                    },
                                                                    "b": {
                                                                        "type": "number"
                                                                    },
                                                                },
                                                                "required": ["a", "b"],
                                                            },
                                                        }
                                                    ]
                                                },
                                            },
                                        },
                                    },
                                }
                            },
                        },
                        "400": {
                            "description": "Bad request - invalid JSON-RPC format",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "jsonrpc": {
                                                "type": "string",
                                                "enum": ["2.0"],
                                            },
                                            "id": {"type": "integer"},
                                            "error": {
                                                "type": "object",
                                                "properties": {
                                                    "code": {"type": "integer"},
                                                    "message": {"type": "string"},
                                                },
                                            },
                                        },
                                    }
                                }
                            },
                        },
                    },
                }
            }
        },
        "components": {
            "schemas": {
                "MCPTool": {
                    "type": "object",
                    "properties": {
                        "name": {"type": "string", "description": "Tool name"},
                        "description": {
                            "type": "string",
                            "description": "Tool description",
                        },
                        "inputSchema": {
                            "type": "object",
                            "description": "JSON Schema for tool input parameters",
                        },
                    },
                }
            }
        },
    }

    return openapi_spec


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
