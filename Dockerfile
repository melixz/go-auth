# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go-auth ./go-auth

WORKDIR /src/go-auth

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseInternal -g cmd/go-auth/main.go

RUN CGO_ENABLED=0 go build -o /go/bin/go-auth ./cmd/go-auth

FROM alpine:latest
WORKDIR /app

COPY --from=builder /go/bin/go-auth ./go-auth

COPY --from=builder /src/go-auth/docs ./docs
COPY --from=builder /src/go-auth/migrations ./migrations

COPY .env .env

EXPOSE 8080 

CMD ["./go-auth"] 