# Makefile for backend_app

# Variables
DOCKER_COMPOSE = docker-compose
SWAG_BIN = ~/go/bin/swag
SWAG_MAIN = cmd/main.go
SWAG_OUT = docs
MIGRATION_DIR = migrations
DB_URL = $(DATABASE_URL) 

.PHONY: all build up down restart logs migrate migrate-up migrate-down migrate-status swagger clean-swagger

all: build

# Docker Compose
build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

restart: down up

logs:
	$(DOCKER_COMPOSE) logs -f

# Database Migration (using golang-migrate)
migrate:
	migrate -path $(MIGRATION_DIR) -database $(DB_URL) up

migrate-up:
	migrate -path $(MIGRATION_DIR) -database $(DB_URL) up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database $(DB_URL) down 1

migrate-status:
	migrate -path $(MIGRATION_DIR) -database $(DB_URL) status

# Swagger Docs
swagger:
	$(SWAG_BIN) init -g $(SWAG_MAIN) -o $(SWAG_OUT)

# Clean up generated docs
clean-swagger:
	rm -rf $(SWAG_OUT)/* 