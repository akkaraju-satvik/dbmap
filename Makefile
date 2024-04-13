build:
	@go build -o bin/dbmap main.go

test:
	@go test ./...