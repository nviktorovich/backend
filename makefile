#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.yaml
LINKS_TABLE ?= crypto_table


#========================#
#== DATABASE MIGRATION ==#
#========================#

migrate-up: ## Run migrations UP
migrate-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

migrate-down: ## Rollback migrations against non test DB
migrate-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down

migrate-create: ## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate create -ext sql -dir /migrations $(LINKS_TABLE)

shell-db: ## Enter to database console
shell-db:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec db psql -U postgres -d postgres

# migrate -path deployment/migrations -database "postgres://postgres:postgres@localhost:5432/links_dev?sslmode=disable" up


# migrate -path deployment/migrations -database "postgres://postgres:postgres@localhost:5432/links_dev?sslmode=disable" up