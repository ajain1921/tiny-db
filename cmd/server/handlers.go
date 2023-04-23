package main

import (
	"fmt"
	"pingack/mp3/internal/config"
	"pingack/mp3/internal/server"
)

type Server struct {
	config   *config.ServerConfigEntry
	database *Database
	bruh     map[string]int
}

func (s *Server) Begin(args *server.BeginArgs, reply *server.Reply) error {
	*reply = "OK"

	s.bruh[args.ClientId] += 1
	fmt.Printf("%+v\n", s.bruh)

	return nil
}

func (s *Server) Deposit(args *server.DepositArgs, reply *server.Reply) error {
	fmt.Println(*args)
	if args.Branch != s.config.Branch {
		// Forward to other server
		// s.servers.Call("Server.Deposit", args, reply)

		return nil
	}

	fmt.Println(*args)

	// s.database.accounts[args.Account]
	return nil
}
