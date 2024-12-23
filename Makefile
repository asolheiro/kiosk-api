include .env

.PHONY: help create-migration migrate-up migrate-down migrate-force

create-migration: ## Create an empty migration
	@read -p "Enter the sequence name: " SEQ; \
    migrate create -ext sql -dir ./internal/pgstore/migrations -seq $${SEQ}

migrate-up:
	@migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" up

migrate-down:
	@read -p "Number of migrations you want to rollback (default: 1): " NUM; NUM=$${NUM:-1}; \
	migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" down $${NUM}

migrate-force:
	@read -p "Enter the version to force: " VERSION; \
	migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" force $${VERSION}

generate:
	sqlc generate -f ./internal/pgstore/sqlc.yaml
