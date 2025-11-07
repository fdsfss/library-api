FROM golang:1.25-alpine AS builder
WORKDIR /library-api

COPY go.mod go.sum ./
RUN go mod download
ENV PORT="8080"
ENV DB_CONN="host=localhost port=5555 user=Dana password=qwerty123 dbname=library-db sslmode=disable"

COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o app ./cmd/

EXPOSE 8080

CMD ["./app"]