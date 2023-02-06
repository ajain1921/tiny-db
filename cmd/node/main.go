package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	if len(args) != 3 {
		return errors.New("exactly 3 command line arguments please")
	}

	name := args[0]
	ip := args[1]
	port := args[2]

	stdin := bufio.NewReader(os.Stdin)

	var timestamp, hash string

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	fmt.Fprintln(conn, name)
	// time.Sleep(4 * time.Second)
	// status, err := bufio.NewReader(conn).ReadString('\n')

	// client, err := rpc.DialHTTP("tcp", ip+":"+port)
	// if err != nil {
	// 	return err
	// }

	// connectData := &logger.ConnectArgs{NodeName: name}
	// client.Go("Logger.Connect", connectData, nil, nil)

	for {
		_, err := fmt.Fscanf(stdin, "%s %s\n", &timestamp, &hash)
		if err != nil {
			fmt.Println("err: ", err)
			break
		}

		// Event := timestamp + " " + hash
		// logData := &logger.LogArgs{Event: Event, NodeName: name}

		// client.Go("Logger.Log", logData, nil, nil)
		fmt.Fprintln(conn, name+" "+timestamp+" "+hash)
	}

	return nil
}
