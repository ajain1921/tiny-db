package main

import "pingack/mp3/internal/server"

// remove index i from array
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func isAbortReply(reply *server.Reply) bool {
	return *reply == "ABORTED" || *reply == "NOT FOUND, ABORTED"
}
