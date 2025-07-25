API_KEY := $(DEEPSEEK_API_KEY)

build:
	go build -o bin/mcp-host ./cmd/main.go

run: build
	./script/run.sh

# Docker rules
docker-build:
	docker build -t mcp-host .

docker-run: docker-build
	docker run -p 9000:9000 -e DEEPSEEK_API_KEY=$(API_KEY) --name mcp-host mcp-host

docker-stop:
	docker stop mcp-host || true
	docker rm mcp-host || true

docker-clean: docker-stop
	docker rmi mcp-host || true

.PHONY: build run docker-build docker-run docker-stop docker-clean