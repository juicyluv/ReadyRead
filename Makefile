build:
	go build -o bin/readyread cmd/main.go

run:
	go run cmd/main.go

test:
	(go test -v -race -timeout 1m -coverprofile cover.out ./internal/...; go tool cover -html=cover.out -o cover.html; rm cover.out)