FROM golang:alpine AS builder

WORKDIR /app

COPY app /app

RUN go mod download
RUN apk --no-cache add bash gcc musl-dev
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN cd /app
RUN go build -o build/gffbot cmd/main.go

FROM alpine AS runner

WORKDIR /gff

COPY --from=builder /app/build/gffbot /gff/build/gffbot
COPY --from=builder /app/configs /gff/configs
COPY --from=builder /app/internal/storage/migrations /gff/migrations
COPY --from=builder /app/makefile /gff

COPY .env ./.env
COPY --from=builder /go/bin/goose /usr/local/bin/goose

RUN apk --no-cache add make
RUN make goose-up

CMD [ "/gff/build/gffbot" ]