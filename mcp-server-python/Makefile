setup:
	./scripts/setup.sh

run: setup
	./scripts/run.sh

# Docker rules
docker-build:
	docker build -t mcp-server-python .

docker-run: docker-build
	docker run -p 9000:9000 -p 9001:9001 --name mcp-server-python mcp-server-python

docker-stop:
	docker stop mcp-server || true
	docker rm mcp-server || true

docker-clean: docker-stop
	docker rmi mcp-server-python || true

.PHONY: setup run docker-build docker-run docker-stop docker-clean