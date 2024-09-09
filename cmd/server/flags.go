package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type netAddress struct {
	host string
	port int
}

func (o *netAddress) String() string {
	return fmt.Sprintf("%s:%d", o.host, o.port)
}

func (o *netAddress) Set(value string) error {
	host, port, err := parseHostAndPort(value)
	if err != nil {
		return err
	}

	o.host = host
	o.port = port

	return nil
}

const (
	defaultHost = "localhost"
	defaultPort = 8080

	envAddress = "ADDRESS"
)

var addr = &netAddress{
	host: defaultHost,
	port: defaultPort,
}

func parseFlags() {
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Server address: host:port")
	flag.Parse()
}

func parseEnvs() {
	if address := os.Getenv(envAddress); address != "" {
		host, port, err := parseHostAndPort(address)
		if err != nil {
			log.Fatal(err)
		}

		addr.host = host
		addr.port = port
	}
}

func parseHostAndPort(s string) (string, int, error) {
	values := strings.Split(s, ":")
	if len(values) != 2 {
		return "", 0, fmt.Errorf("invalid format")
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return "", 0, err
	}

	return values[0], port, nil
}
