all:
	go mod tidy
	go mod verify
	go fmt ./...
	go test -cover ./...
	go install ./...

lint:
	go mod tidy
	go mod verify
	golangci-lint run

test:
	go mod tidy
	go test -v --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
