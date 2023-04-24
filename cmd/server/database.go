package main

import (
	"fmt"
	"sync"
)

type Database struct {
	lock     sync.Mutex
	accounts map[string]*Account
}

func (d *Database) printDatabase() {
	for accountName, account := range d.accounts {
		fmt.Printf("%s: {balance: %d, creators: %v, RTS: %v, TW: %v}\n", accountName, account.balance, account.creators, account.readTimestamps, account.tentativeWrites)
	}
	fmt.Printf("\n")
}

type TentativeWrite struct {
	timestamp        int64
	tentativeBalance int
}

type Account struct {
	lock    sync.Mutex
	balance int
	name    string

	creators Set[int64]

	committedBalance   int
	committedTimestamp int64
	readTimestamps     Set[int64]
	tentativeWrites    []*TentativeWrite
}

type Set[T comparable] map[T]interface{}

func (s *Set[T]) Add(key T) {
	(*s)[key] = struct{}{}
}

func (s *Set[T]) Remove(key T) {
	delete(*s, key)
}

func (s *Set[T]) Has(key T) bool {
	_, exists := (*s)[key]
	return exists
}
