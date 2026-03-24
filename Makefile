.PHONY: build test vet lint sec check docker-up docker-build clean

## Build
build:
	go build ./cmd/pedmin/...

## Test
test:
	go test ./...

## Static analysis
vet:
	go vet ./...

lint:
	golangci-lint run

sec:
	gosec ./...

## Run all checks (vet + lint + sec)
check: vet lint sec

## Docker
docker-up:
	docker compose up

docker-build:
	docker compose build

## Clean build cache
clean:
	go clean -cache
