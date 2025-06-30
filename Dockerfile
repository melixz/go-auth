# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency --parseInternal -g cmd/go-auth/main.go || true
RUN go build -o go-auth ./cmd/go-auth

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/go-auth ./go-auth
COPY --from=builder /app/docs ./docs
COPY .env .env
EXPOSE 8080
CMD ["./go-auth"] 