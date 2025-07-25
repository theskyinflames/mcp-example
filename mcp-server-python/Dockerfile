FROM python:3.11-bookworm AS builder

# Install make in builder stage
RUN apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

# Set working directory for builder stage
WORKDIR /app

# Install dependencies in builder stage
COPY . .
RUN pip install --no-cache-dir --user -r requirements.txt

# Production stage with distroless image
FROM gcr.io/distroless/python3-debian12:latest

# Set working directory
WORKDIR /app

# Copy Python packages from builder stage
COPY --from=builder /root/.local /root/.local

# Copy the rest of the application
COPY . .

# Expose port 9000 (default for uvicorn)
EXPOSE 9000
EXPOSE 9001

# Update PATH to include user packages
ENV PATH=/root/.local/bin:$PATH
ENV PYTHONPATH=/root/.local/lib/python3.11/site-packages

CMD ["mcp-server.py"]