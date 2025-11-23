APP_NAME := pr-reviewer-service
CMD_DIR  := ./cmd/pr-reviewer-service

.PHONY: build run test swag docker-build docker-up docker-down docker-logs

build:
	go build -v -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)/main.go

swag:
	swag init -g cmd/pr-reviewer-service/main.go -o ./docs

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

docker-logs:
	docker compose logs -f app

.DEFAULT_GOAL := run