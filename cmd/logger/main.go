package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
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

	if len(args) != 1 {
		return errors.New("exactly 1 command line argument please")
	}

	port := args[0]

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ruh roh: %v\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	var nodeName string
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, " ", 2)

		nodeName = split[0]

		if len(split) == 1 {
			printConnectionChange(nodeName, "connected")
			continue
		}

		splitMessage := strings.Split(split[1], " ")

		fmt.Println(splitMessage[0], nodeName, splitMessage[1])
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	printConnectionChange(nodeName, "disconnected")
}

func printConnectionChange(nodeName, connectionType string) {
	now := time.Now()
	timestamp := fmt.Sprintf("%f", float64(now.UnixNano())/math.Pow(10, 9))
	fmt.Println(timestamp, "-", nodeName, connectionType)
}
