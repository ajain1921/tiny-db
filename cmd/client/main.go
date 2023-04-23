package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"pingack/mp3/internal/config"
	"pingack/mp3/internal/server"
	"strings"
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

	// coordinatingServer := servers[rand.Intn(len(servers))]
	coordinatingServer := servers[0]

	fmt.Println(coordinatingServer)

	client, err := rpc.DialHTTP("tcp", coordinatingServer.Hostname+":"+coordinatingServer.Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// var server *config.Server
	var input string
	reader := bufio.NewReader(os.Stdin)
	reader.Reset(os.Stdin)
	for {
		input, err = reader.ReadString('\n')
		input = input[:len(input)-1]
		if err != nil {
			return err
		}

		switch command := strings.Split(input, " ")[0]; command {
		case "BEGIN":
			args := server.BeginArgs{ClientId: id}
			var reply server.Reply
			err = client.Call("Server.Begin", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

		case "DEPOSIT":
			var any string
			var branchAndAccount string
			var amount int

			fmt.Sscanf(input, "%s %s %d", &any, &branchAndAccount, &amount)

			branch := strings.Split(branchAndAccount, ".")[0]
			account := strings.Split(branchAndAccount, ".")[1]

			args := server.DepositArgs{Branch: branch, Account: account, Amount: amount}
			var reply server.Reply
			err = client.Call("Server.Deposit", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}

		case "COMMIT":
			/* if server == nil {
				continue
			} */
			return nil

		case "":

		}
	}

	return nil
}
