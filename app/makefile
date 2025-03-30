.SILENT:

#include .env

go-run:
	go run cmd/main.go

build:
	go build -o build/gffbot cmd/main.go

run:
	./build/gffbot

goose-up:
	goose -dir ./migrations postgres 'user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) port=5432 host=db sslmode=disable' up

goose-down:
	goose -dir ./migrations postgres 'user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) port=5432 host=db sslmode=disable' down

test:
	go test ./...
