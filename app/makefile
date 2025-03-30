.SILENT:

#include .env

go-run:
	go run cmd/main.go

build:
	go build -o build/gffbot cmd/main.go

run:
	./build/gffbot

goose-up:
	goose -dir ./migrations postgres 'user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) port=$(POSTGRES_PORT) host=$(DB_HOST) sslmode=disable' up

goose-down:
	goose -dir ./migrations postgres 'user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) port=$(POSTGRES_PORT) host=$(DB_HOST) sslmode=disable' down

test:
	go test ./...
