.PHONY: build run

build:
	go build -o ./build/chat ./cmd/chat

run:
	make build && ./build/chat -addr :8080