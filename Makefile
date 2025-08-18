API_KEY := $(DEEPSEEK_API_KEY)

run-demo:
	scripts/run-demo-interactive.sh

down:
	docker-compose down

clean:
	docker-compose down --rmi all --volumes --remove-orphans

go-server-logs:
	docker-compose logs -f mcp-server-go

python-server-logs:
	docker-compose logs -f mcp-server-python

host-logs:
	docker-compose logs -f mcp-host

.PHONY: up down clean

