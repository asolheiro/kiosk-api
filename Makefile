migrate:
	@migrate -path=./internal/pgstore/migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

generate:
	@sqlc generate -f ./internal/pgstore/sqlc.yaml
