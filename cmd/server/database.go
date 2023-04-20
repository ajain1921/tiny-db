package main

import "sync"

type Database struct {
	accounts map[string]*Account
}

type Account struct {
	lock    sync.RWMutex
	balance int
	name    string
}
