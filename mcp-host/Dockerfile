# ---- Builder Stage ----
# Use the official Golang image to create a build environment
FROM golang:1.24.5-alpine3.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application, creating a statically linked binary
# CGO_ENABLED=0 is important for running in a scratch container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/mcp-host ./cmd/main.go


# ---- Production Stage ----
# Use a minimal scratch image for the final container
FROM scratch

# Set the working directory
WORKDIR /

# Copy the CA certificates from the builder stage
# This is necessary for making HTTPS requests (e.g., to the LLM API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the built binary from the builder stage
COPY --from=builder /bin/mcp-host /host-app

# Set the entrypoint for the container
ENTRYPOINT ["/host-app"]
