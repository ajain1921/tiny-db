package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"pingack/mp3/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	if len(args) != 2 {
		return errors.New("exactly 2 command line arguments please")
	}

	branch := args[0]
	configFilename := args[1]

	servers, err := config.ParseConfig(configFilename)
	if err != nil {
		return err
	}

	var serverConfigEntry *config.ServerConfigEntry
	for _, potentialServer := range servers {
		if potentialServer.Branch == branch {
			serverConfigEntry = potentialServer
		}
	}

	server := &Server{config: serverConfigEntry}
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", server.config.Hostname+":"+server.config.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

	return nil
}
