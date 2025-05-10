BINARY_NAME = server
DOCKER_IMAGE = calc-server
PORT = 8080
NETWORK = calc-net

.PHONY: all run build test clean docker_build docker_run docker_clean help

start: run

run:
	go run -v ./cmd/server

build:
	go build -v -o $(BINARY_NAME) ./cmd/server

test: unit-test integration-test

unit-test:
	go test -v ./server/internal/... -coverprofile=coverage.out

integration-test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

coverage:
	go tool cover -html=coverage.out

swagger:
	swag init -g cmd/server/main.go --output docs

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker_clean
	docker run --name $(BINARY_NAME) \
		-p $(PORT):$(PORT) \
		--network $(NETWORK) \
		-e TIME_ADDITION_MS=1000 \
		-e TIME_SUBTRACTION_MS=1000 \
		-e TIME_MULTIPLICATION_MS=1000 \
		-e TIME_DIVISION_MS=1000 \
		-d $(DOCKER_IMAGE)

docker-clean:
	docker rm -f $(BINARY_NAME) 2>/dev/null || true

docker-network:
	docker network create $(NETWORK) 2>/dev/null || true

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)