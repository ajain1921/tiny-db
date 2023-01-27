all: node logger
.PHONY: all

node: cmd/node/main.go node.go
	go build -o node cmd/node/main.go

logger: cmd/logger/main.go logger.go
	go build -o logger cmd/logger/main.go
