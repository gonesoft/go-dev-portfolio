# Variables
DOCKER_COMPOSE=docker compose

# DB up
db-up:
	$(DOCKER_COMPOSE) up -d

# db down
db-down:
	$(DOCKER_COMPOSE) down

# run app
run:
	go run ./cmd/api

# Instalar dependencias
deps:
	go mod tidy

# reset db, this will erase all data
db-reset:
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up -d

# run db test
db-test-up:
	$(DOCKER_COMPOSE) up -d postgres-test

# run db test down
db-test-down:
	$(DOCKER_COMPOSE) down -v postgres-test

db-up:
	$(DOCKER_COMPOSE) up -d

test:
	go test -v ./...