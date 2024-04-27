all:
	go mod tidy
	go mod verify
	go fmt ./...
	go vet ./...
	go test -cover ./...
	go install ./...

lint:
	go mod tidy
	go mod verify
	golangci-lint run
	gocyclo -over 15 .

test:
	go mod tidy
	go test -v --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
