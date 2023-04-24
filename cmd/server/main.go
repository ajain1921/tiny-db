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
	"time"
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
			continue
		}
	}

	connections := make(map[string]*rpc.Client)
	server := &Server{config: serverConfigEntry, servers: connections, transactions: make(map[int64](Set[string])), database: &Database{accounts: make(map[string]*Account)}}

	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", server.config.Hostname+":"+server.config.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	time.Sleep(time.Duration(3) * time.Second)

	for _, potentialServer := range servers {
		client, err := rpc.DialHTTP("tcp", potentialServer.Hostname+":"+potentialServer.Port)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		connections[potentialServer.Branch] = client
	}

	fmt.Println("Ready")
	for {

	}

}
