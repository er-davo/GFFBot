FROM golang:alpine

WORKDIR /app

COPY gff ./

RUN go mod download
RUN apk --no-cache add make bash gcc musl-dev

COPY GFF-project ./

CMD [ "make", "run" ]