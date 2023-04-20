all: node logger
.PHONY: all

client: cmd/client/main.go internal
	go build -o bin/client cmd/client/main.go

server: cmd/server/main.go internal
	go build -o bin/server cmd/server/main.go

internal: internal/config/config.go
