all: node logger
.PHONY: all

node: node.go
	go build node.go

logger: logger.go
	go build logger.go
