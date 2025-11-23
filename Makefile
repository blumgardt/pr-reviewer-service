.PHONY: build
build:
		go build -v ./cmd/pr-reviewer-service
run:
		go run ./cmd/pr-reviewer-service/main.go

.DEFAULT_GOAL := run