-include .env

COMPOSE        = docker compose --env-file .env -f deploy/compose.yaml
MIGRATIONS_DIR = services/core/migrations
DB_URL         = postgres://$(CORE_POSTGRES_USER):$(CORE_POSTGRES_PASSWORD)@localhost:$(CORE_POSTGRES_PORT)/$(CORE_POSTGRES_DB)?sslmode=disable

.PHONY: up down reset logs ps run test psql \
        migrate-create migrate-up migrate-down

up: 
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

reset:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

run: 
	go run ./services/core/cmd/core

# test:          
# 	go test -race ./...

psql:
	$(COMPOSE) exec postgres psql -U $(CORE_POSTGRES_USER) -d $(CORE_POSTGRES_DB)

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=create_users" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(CORE_DATABASE_URL)" up

migrate-down: 
	migrate -path $(MIGRATIONS_DIR) -database "$(CORE_DATABASE_URL)" down 1