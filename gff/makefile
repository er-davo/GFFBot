.SILENT:

go-run:
	go run cmd/main.go

build:
	go build -o build/gffbot cmd/main.go

run:
	./build/gffbot

goose-up:
	goose postgres 'user=user password=password dbname=gff sslmode=disable' up

goose-down:
	goose -dir ./migrations postgres 'user=user password=password dbname=gff sslmode=disable' down

test:
	go test ./...
