API_KEY := $(DEEPSEEK_API_KEY)

run-demo:
	scripts/run-demo.sh

down:
	docker-compose down

clean:
	docker-compose down --rmi all --volumes --remove-orphans

server-logs:
	docker-compose logs -f mcp-server

host-logs:
	docker-compose logs -f mcp-host

.PHONY: up down clean

