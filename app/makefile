.SILENT:

#include .env

go-run:
	go run cmd/main.go

build:
	go build -o build/gffbot cmd/main.go

run:
	./build/gffbot

goose-up1:
	goose -dir ./migrations postgres "$$DATABASE_URL" up

goose-up:
	@test -n "$$DATABASE_URL" || (echo "Error: DATABASE_URL is not set"; exit 1)
    goose postgres "$$DATABASE_URL" up

goose-down:
	goose -dir ./migrations postgres 'user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB port=$$POSTGRES_PORT host=$$DB_HOST sslmode=disable' down

test:
	go test ./...
