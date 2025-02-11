include .env

.PHONY: help create-migration migrate-up migrate-down migrate-force

create-migration: ## Create an empty migration
	@read -p "Enter the sequence name: " SEQ; \
    migrate create -ext sql -dir ./internal/pgstore/migrations -seq $${SEQ}

migrate-up:
	@migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" up

migrate-up-sqlite:
	@migrate -path=./internal/sqlitestore/migrations -database "sqlite://${SQLITE_DB_PATH}" up


migrate-down:
	@read -p "Number of migrations you want to rollback (default: 1): " NUM; NUM=$${NUM:-1}; \
	migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" down $${NUM}

migrate-force:
	@read -p "Enter the version to force: " VERSION; \
	migrate -path=./internal/pgstore/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL}" force $${VERSION}

generate:
	sqlc generate -f ./internal/pgstore/sqlc.yaml

generate-sqlite:
	sqlc generate -f ./internal/sqlitestore/sqlc.yaml


run: 
	migrate-up-sqlite
	@clear
	@go run cmd/kiosk/kiosk.go

swagger:
	docker run --rm -v $(shell pwd):/code ghcr.io/swaggo/swag:latest init -g ./cmd/kiosk/kiosk.go
	rm -f ./docs/docs.json && rm -f ./docs/docs.yaml
	mv ./docs/swagger.json ./docs/docs.json && mv ./docs/swagger.yaml ./docs/docs.yaml
	sed -i '' 's|"github.com/swaggo/swag/v2"|"github.com/swaggo/swag"|' ./docs/docs.go


