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
	go run main.go

# Instalar dependencias
deps:
	go mod tidy

# reset db, this will erase all data
db-reset:
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up -d