build:
	go build -o bin/mcp-server ./cmd/main.go

run: build
	./script/run.sh

# Docker rules
docker-build:
	docker build -t mcp-server .

docker-run: docker-build
	docker run -p 8090:8090 -p 8091:8091 -e DEEPSEEK_API_KEY=$(API_KEY) --name mcp-server mcp-server

docker-stop:
	docker stop mcp-server || true
	docker rm mcp-server || true

docker-clean: docker-stop
	docker rmi mcp-server || true

.PHONY: build run docker-build docker-run docker-stop docker-clean