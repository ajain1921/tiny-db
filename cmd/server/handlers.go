package main

import (
	"fmt"
	"log"
	"net/rpc"
	"pingack/mp3/internal/config"
	"pingack/mp3/internal/server"
	"sync"
	"time"
)

type Server struct {
	config           *config.ServerConfigEntry
	database         *Database
	servers          map[string]*rpc.Client
	transactionsLock sync.Mutex
	transactions     map[int64](Set[string]) // client to branches that it's interacted with
}

func (s *Server) Begin(args *server.BeginArgs, reply *server.Reply) error {
	s.transactionsLock.Lock()
	defer s.transactionsLock.Unlock()
	*reply = "OK"

	s.transactions[args.Timestamp] = make(Set[string])

	return nil
}

func (s *Server) Deposit(args *server.UpdateArgs, reply *server.Reply) error {
	err := s.update(args, reply, true)
	s.detectAbort(args.Timestamp, reply)

	s.database.printDatabase()
	return err
}

func (s *Server) Withdraw(args *server.UpdateArgs, reply *server.Reply) error {
	err := s.update(args, reply, false)
	s.detectAbort(args.Timestamp, reply)
	return err
}

func (s *Server) Balance(args *server.BalanceArgs, reply *server.Reply) error {
	fmt.Println(*args)

	if args.Branch != s.config.Branch {
		// Forward to other server
		s.servers[args.Branch].Call("Server.Balance", args, reply)
		return nil
	}

	//deposit logic
	if _, ok := s.database.accounts[args.Account]; !ok {
		// if account doesn't exist already, ABORT!
		*reply = "NOT FOUND, ABORTED"
		s.detectAbort(args.Timestamp, reply)
		return nil
	}

	account := s.database.accounts[args.Account]

	fmt.Println("Account found", account)

	for {
		account.lock.Lock()

		if args.Timestamp > account.committedTimestamp {

			maxViableWrite := &TentativeWrite{timestamp: account.committedTimestamp, tentativeBalance: account.committedBalance}

			for _, tentativeWrite := range account.tentativeWrites {
				if tentativeWrite.timestamp > maxViableWrite.timestamp && tentativeWrite.timestamp <= args.Timestamp {
					maxViableWrite = tentativeWrite
				}
			}

			fmt.Println("Max viable write timestamp", maxViableWrite.timestamp)

			if maxViableWrite.timestamp == account.committedTimestamp {
				*reply = server.Reply(fmt.Sprintf("%s.%s = %d", s.config.Branch, account.name, maxViableWrite.tentativeBalance))
				account.readTimestamps.Add(args.Timestamp)

				account.lock.Unlock()
				return nil
			} else {
				fmt.Println("Why no equal?", maxViableWrite.timestamp, args.Timestamp)
				if maxViableWrite.timestamp == args.Timestamp {
					*reply = server.Reply(fmt.Sprintf("%s.%s = %d", s.config.Branch, account.name, maxViableWrite.tentativeBalance))
					account.lock.Unlock()
					return nil
				} else {
					// wait until the transaction that wrote Ds is committed or aborted, and
					// reapply the read rule.
					// if the transaction is committed, Tc will read its value after the wait.
					// if the transaction is aborted, Tc will read the value from an older
					// transaction.
				}
			}
		} else {
			*reply = "ABORTED"
			account.lock.Unlock()

			s.detectAbort(args.Timestamp, reply)
			return nil
		}

		account.lock.Unlock()
		time.Sleep(time.Second * time.Duration(1))
	}
}

func (s *Server) Abort(args *server.AbortArgs, reply *server.Reply) error {
	s.abort(args.Timestamp)
	*reply = "ABORTED"
	return nil
}

func (s *Server) update(args *server.UpdateArgs, reply *server.Reply, deposit bool) error {
	fmt.Println(*args)

	s.transactionsLock.Lock()

	if _, ok := s.transactions[args.Timestamp]; ok {
		// we are coordinator so add branch to set
		branchSet := s.transactions[args.Timestamp]
		branchSet.Add(args.Branch)
	}

	s.transactionsLock.Unlock()

	if args.Branch != s.config.Branch {
		function := "Server.Deposit"
		if !deposit {
			function = "Server.Withdraw"
		}
		// Forward to other server
		s.servers[args.Branch].Call(function, args, reply)
		return nil
	}

	//deposit logic
	if _, ok := s.database.accounts[args.Account]; !ok {
		// if account doesn't exist already, initialize
		s.database.accounts[args.Account] = &Account{name: args.Account, creators: Set[int64]{args.Timestamp: struct{}{}}, readTimestamps: make(Set[int64])}
	}
	account := s.database.accounts[args.Account]
	account.lock.Lock()
	defer account.lock.Unlock()

	maxReadTimestamp := int64(-1)
	for timestamp := range account.readTimestamps {
		if timestamp > maxReadTimestamp {
			maxReadTimestamp = timestamp
		}
	}

	if args.Timestamp >= maxReadTimestamp && args.Timestamp > account.committedTimestamp {
		*reply = "OK"

		amount := args.Amount
		if !deposit {
			amount *= -1
		}

		for _, tentativeWrite := range account.tentativeWrites {
			if tentativeWrite.timestamp == args.Timestamp {
				tentativeWrite.tentativeBalance += amount
				return nil
			}
		}

		tentativeWrite := &TentativeWrite{
			tentativeBalance: account.committedBalance + amount,
			timestamp:        args.Timestamp,
		}
		account.tentativeWrites = append(account.tentativeWrites, tentativeWrite)
		if account.committedTimestamp == 0 {
			account.creators.Add(args.Timestamp)
		}
		return nil
	}

	*reply = "ABORTED"
	return nil
}

// for coordinators to handle abort
func (s *Server) detectAbort(timestamp int64, reply *server.Reply) {
	if _, ok := s.transactions[timestamp]; !ok {
		// we're not coordinator
		return
	}

	if *reply == "ABORTED" || *reply == "NOT FOUND, ABORTED" {
		s.transactionsLock.Lock()

		for branch := range s.transactions[timestamp] {
			if branch == s.config.Branch {
				continue
			}

			args := server.AbortArgs{Timestamp: timestamp}
			err := s.servers[branch].Call("Server.Abort", args, reply)
			if err != nil {
				log.Fatal("ERROR!")
			}
		}

		s.transactionsLock.Unlock()

		s.abort(timestamp)
	}
}

func (s *Server) abort(timestamp int64) {
	s.database.lock.Lock()
	defer s.database.lock.Unlock()
	for accountName, account := range s.database.accounts {

		account.lock.Lock()

		account.creators.Remove(timestamp)
		if len(account.creators) == 0 && account.committedTimestamp == 0 {
			delete(s.database.accounts, accountName)
			account.lock.Unlock()
			continue
		}

		// remove from readTimestamps
		account.readTimestamps.Remove(timestamp)

		// remove from tentativeWrites
		tentativeWriteIdx := -1
		for i, tentativeWrite := range account.tentativeWrites {
			if timestamp == tentativeWrite.timestamp {
				tentativeWriteIdx = i
				break
			}
		}

		if tentativeWriteIdx != -1 {
			account.tentativeWrites = remove(account.tentativeWrites, tentativeWriteIdx)
		}

		account.lock.Unlock()
	}

	s.transactionsLock.Lock()
	delete(s.transactions, timestamp)
	s.transactionsLock.Unlock()
}

// remove index i from array
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
