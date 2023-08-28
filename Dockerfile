FROM golang:latest as builder
WORKDIR /service

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o service ./cmd/service/service.go

FROM ubuntu AS production

WORKDIR /service

RUN apt-get update

COPY --from=builder /service/service ./
COPY --from=builder /service/config.env ./

CMD ["./service"]