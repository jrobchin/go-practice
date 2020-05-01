.PHONY: build run

build:
	cd src && go build -o ../build/chat

run:
	make build && cd src && ../build/chat -addr :8080