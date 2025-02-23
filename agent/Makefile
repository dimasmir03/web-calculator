BINARY_NAME = agent
DOCKER_IMAGE = calc-agent
SERVER_URL = http://calc-server:8080
NETWORK = calc-net

.PHONY: all run build test clean docker_build docker_run docker_clean help

all: run

run:
	go run -v ./cmd/agent

build:
	go build -v -o $(BINARY_NAME) ./cmd/agent
	./server

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

clean:
	rm -rf $(BINARY_NAME) coverage.out

docker_build:
	docker build -t $(DOCKER_IMAGE) .

docker_run: docker_clean
	docker run --name $(BINARY_NAME) \
 		--network $(NETWORK) \
 		-e SERVER_URL=$(SERVER_URL) \
 		-d $(DOCKER_IMAGE)

docker_clean:
	docker rm -f $(BINARY_NAME) 2>/dev/null || true

docker_network:
	docker network create $(NETWORK) 2>/dev/null || true

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)