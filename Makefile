.PHONY: run test cover fmt vet lint build clean

run:
	go run ./cmd/app

test:
	go test ./...

cover:
	go test ./... -cover

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -o bin/metadata-ingestion-service ./cmd/app

clean:
	rm -rf bin coverage.out

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out