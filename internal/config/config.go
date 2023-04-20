package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ServerConfigEntry struct {
	Branch   string
	Hostname string
	Port     string
}

func ParseConfig(filename string) ([]*ServerConfigEntry, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	servers := make([]*ServerConfigEntry, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()

		var branch string
		var hostname string
		var port string
		_, err = fmt.Fscanf(strings.NewReader(line), "%s %s %s", &branch, &hostname, &port)
		if err != nil {
			return nil, err
		}

		servers = append(servers, &ServerConfigEntry{Branch: branch, Hostname: hostname, Port: port})
	}
	return servers, nil
}
