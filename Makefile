run:
	go run .

lint-install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	@golangci-lint run ./...
