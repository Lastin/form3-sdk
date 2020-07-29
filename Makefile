.PHONY: test

start-stack:
	docker-compose up -d

test: start-stack
	@go test ./...