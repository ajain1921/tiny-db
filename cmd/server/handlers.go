package main

import (
	"pingack/mp3/internal/config"
	"pingack/mp3/internal/server"
)

type Server struct {
	config   *config.ServerConfigEntry
	database *Database
}

func (s *Server) Begin(args *server.BeginArgs, reply *server.Reply) error {
	*reply = "OK"
	return nil
}

func (s *Server) Deposit(args *server.DepositArgs, reply *server.Reply) error {
	if args.Branch != s.config.Branch {
		// Forward to other server
		// s.servers.Call("Server.Deposit", args, reply)

		return nil
	}

	s.database.accounts[args.Account]
	return nil
}
