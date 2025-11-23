FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/pr-reviewer-service/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-reviewer-service ./cmd/pr-reviewer-service



FROM alpine:3.20

RUN apk add --no-cache ca-certificates curl postgresql-client netcat-openbsd

WORKDIR /app

COPY --from=builder /app/pr-reviewer-service .
COPY --from=builder /app/docs ./docs
COPY internal/migrations ./internal/migrations
COPY config.toml ./config.toml
COPY docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ENV MIGRATIONS_PATH=/app/internal/migrations

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]