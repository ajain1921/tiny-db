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

const FILENAME = "data.csv"
const HEADER = "timestamp,sent_timestamp,bytes_received\n"

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

	f, err := os.Create("data.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ruh roh: %v\n", err)
	}
	f.WriteString(HEADER)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ruh roh: %v\n", err)
		}
		go handleConnection(conn, f)
	}
}

func handleConnection(conn net.Conn, f *os.File) {
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
		sentTimestamp := splitMessage[0]
		hash := splitMessage[1]

		fmt.Println(sentTimestamp, nodeName, hash)

		numBytes := len([]byte(line))
		writeLine(f, sentTimestamp, numBytes)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	printConnectionChange(nodeName, "disconnected")
}

func printConnectionChange(nodeName, connectionType string) {
	timestamp := fmt.Sprintf("%f", currentTime())
	fmt.Println(timestamp, "-", nodeName, connectionType)
}

func currentTime() float64 {
	now := time.Now()
	return float64(now.UnixNano()) / math.Pow(10, 9)
}

func writeLine(f *os.File, sentTimestamp string, bytesReceived int) {
	s := fmt.Sprintf("%f,%s,%d\n", currentTime(), sentTimestamp, bytesReceived)
	f.WriteString(s)
}
