all: node logger
.PHONY: all

node: cmd/node/main.go
	go build -o bin/node cmd/node/main.go

logger: cmd/logger/main.go
	go build -o bin/logger cmd/logger/main.go
