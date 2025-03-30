.SILENT:

#include .env

go-run:
	go run cmd/main.go

build:
	go build -o build/gffbot cmd/main.go

run:
	./build/gffbot

goose-up:
	goose -dir ./migrations postgres $(DATABASE_URL) up

goose-down:
	goose -dir ./migrations postgres $(DATABASE_URL) down

test:
	go test ./...
