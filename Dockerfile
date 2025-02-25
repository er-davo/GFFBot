FROM golang:alpine AS builder

WORKDIR /app

COPY gff ./

RUN go mod download
RUN apk --no-cache add bash gcc musl-dev
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD [ "go" "build" "-o" "app/build/gffbot" "app/cmd/main.go" ]

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/storage/migrations ./
COPY --from=builder /app/makefile ./

RUN apk --no-cache add make
RUN make goose-up

CMD [ "./build/gffbot" ]