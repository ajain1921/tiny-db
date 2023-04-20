package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"pingack/mp3/internal/config"
	"pingack/mp3/internal/server"
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
		return errors.New("exactly 2 command line argument please")
	}

	id := args[0]
	configFilename := args[1]

	servers, err := config.ParseConfig(configFilename)
	if err != nil {
		return err
	}

	coordinatingServer := servers[rand.Intn(len(servers))]

	client, err := rpc.DialHTTP("tcp", coordinatingServer.Hostname+coordinatingServer.Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// var server *config.Server
	var command string
	stdin := bufio.NewReader(os.Stdin)
	stdin.Reset(os.Stdin)
	for {
		_, err := fmt.Fscanf(stdin, "%s\n", &command)
		if err != nil {
			return err
		}

		switch command {
		case "BEGIN":
			args := server.BeginArgs{}
			var reply server.Reply
			err = client.Call("Server.Begin", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}

		case "DEPOSIT":
			var branch string
			var account string
			var amount int

			fmt.Scanf("%s %s.%s %d", &branch, &account, &amount)
			args := &server.DepositArgs{Branch: branch, Account: account, Amount: amount}
			var reply server.Reply
			err = client.Call("Server.Deposit", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}

		case "COMMIT":
			if server == nil {
				continue
			}
			return nil

		case "":

		}
	}

	return nil
}
