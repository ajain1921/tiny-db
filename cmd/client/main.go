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
	"strconv"
	"strings"
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
		return errors.New("exactly 2 command line argument please")
	}
	rand.Seed(time.Now().UnixNano())

	id := args[0]
	configFilename := args[1]

	servers, err := config.ParseConfig(configFilename)
	if err != nil {
		return err
	}

	coordinatingServer := servers[rand.Intn(len(servers))]
	// coordinatingServer := servers[0]

	client, err := rpc.DialHTTP("tcp", coordinatingServer.Hostname+":"+coordinatingServer.Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	timestamp := time.Now().UnixNano()
	// fmt.Println(coordinatingServer, timestamp, len(servers), rand.Intn(len(servers)))

	// var server *config.Server
	var input string
	reader := bufio.NewReader(os.Stdin)
	reader.Reset(os.Stdin)
	for {
		input, err = reader.ReadString('\n')
		if len(input) == 0 {
			continue
		}

		input = input[:len(input)-1]
		if err != nil {
			return err
		}

		command := strings.Split(input, " ")[0]

		if command == "sleep" {
			seconds, _ := strconv.Atoi(strings.Split(input, " ")[1])
			time.Sleep(time.Second * time.Duration(seconds))
		}

		switch command {
		case "BEGIN":
			args := server.BeginArgs{ClientId: id, Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.Begin", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			if isAborted(reply) {
				return nil
			}

		case "DEPOSIT":
			var any string
			var branchAndAccount string
			var amount int

			fmt.Sscanf(input, "%s %s %d", &any, &branchAndAccount, &amount)

			branch := strings.Split(branchAndAccount, ".")[0]
			account := strings.Split(branchAndAccount, ".")[1]

			args := server.UpdateArgs{Branch: branch, Account: account, Amount: amount, Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.Deposit", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			if isAborted(reply) {
				return nil
			}

		case "WITHDRAW":
			var any string
			var branchAndAccount string
			var amount int

			fmt.Sscanf(input, "%s %s %d", &any, &branchAndAccount, &amount)

			branch := strings.Split(branchAndAccount, ".")[0]
			account := strings.Split(branchAndAccount, ".")[1]

			args := server.UpdateArgs{Branch: branch, Account: account, Amount: amount, Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.Withdraw", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			if isAborted(reply) {
				return nil
			}

		case "BALANCE":
			var any string
			var branchAndAccount string

			fmt.Sscanf(input, "%s %s %d", &any, &branchAndAccount)

			branch := strings.Split(branchAndAccount, ".")[0]
			account := strings.Split(branchAndAccount, ".")[1]

			args := server.BalanceArgs{Branch: branch, Account: account, Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.Balance", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			if isAborted(reply) {
				return nil
			}

		case "ABORT":
			args := server.AbortArgs{Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.Abort", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			return nil

		case "COMMIT":
			args := server.CommitArgs{Timestamp: timestamp}
			var reply server.Reply
			err = client.Call("Server.CoordinateCommit", &args, &reply)
			if err != nil {
				log.Fatal("begin error:", err)
			}
			fmt.Println(reply)

			return nil

		case "":

		}
	}

	return nil
}

func isAborted(reply server.Reply) bool {
	return reply == "ABORTED" || reply == "NOT FOUND, ABORTED"
}
