APP_NAME := quiz-realtime

.PHONY: build run test docker-build docker-run

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run --rm -p 8080:8080 -v $(PWD)/configs:/app/configs $(APP_NAME):latest

