# MCP example with Go

## Get tools list from users MCP server

```sh
curl --request POST \
  --url http://localhost:8090/mcp \
  --header 'Accept: application/json, text/event-stream' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.2.0' \
  --data '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}'
```

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "annotations": {
          "readOnlyHint": false,
          "destructiveHint": true,
          "idempotentHint": false,
          "openWorldHint": true
        },
        "description": "Create a new user",
        "inputSchema": {
          "properties": {
            "age": {
              "minimum": 0,
              "type": "number"
            },
            "email": {
              "type": "string"
            },
            "name": {
              "type": "string"
            },
            "user_id": {
              "type": "string"
            }
          },
          "required": ["user_id", "name", "email"],
          "type": "object"
        },
        "name": "create_user"
      },
      {
        "annotations": {
          "readOnlyHint": false,
          "destructiveHint": true,
          "idempotentHint": false,
          "openWorldHint": true
        },
        "description": "Get user information",
        "inputSchema": {
          "properties": {
            "user_id": {
              "type": "string"
            }
          },
          "required": ["user_id"],
          "type": "object"
        },
        "name": "get_user"
      }
    ]
  }
}
```
