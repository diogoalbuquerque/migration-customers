.PHONY: format build clean test

format:
	go mod tidy
	go fmt

build: format
	export GO111MODULE=on
	GOOS=linux GOARCH=amd64 go build -o bin/MIGRATION_CUSTOMERS ./

clean:
	rm -rf ./.scannerwork ./bin ./vendor Gopkg.lock coverage.out

test: clean format
	go test -v -covermode=count -coverprofile=coverage.out ./...
