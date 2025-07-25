# ---- Builder Stage ----
# Use the official Golang image to create a build environment
FROM golang:1.24.5-alpine3.22 AS builder

# 1. Install netcat in the builder stage, which has a package manager
RUN apk add --no-cache ca-certificates netcat-openbsd

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/mcp-server ./cmd/main.go


# ---- Production Stage ----
# Use a minimal scratch image for the final container
FROM scratch

# Set the working directory
WORKDIR /

# Copy the CA certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 2. Copy the netcat binary from the builder into the final image
COPY --from=builder /usr/bin/nc /bin/nc

# Copy the built application binary from the builder stage
COPY --from=builder /bin/mcp-server /server-app

# Expose the port the server listens on
EXPOSE 8090
EXPOSE 8091

# Set the entrypoint for the container
ENTRYPOINT ["/server-app"]