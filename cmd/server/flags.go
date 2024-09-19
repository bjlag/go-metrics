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
	defaultHost     = "localhost"
	defaultPort     = 8080
	defaultLogLevel = "info"

	envAddress  = "ADDRESS"
	envLogLevel = "LOG_LEVEL"
)

var (
	addr = &netAddress{
		host: defaultHost,
		port: defaultPort,
	}

	logLevel string
)

func parseFlags() {
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Server address: host:port")
	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.Parse()
}

func parseEnvs() {
	if envAddressValue := os.Getenv(envAddress); envAddressValue != "" {
		host, port, err := parseHostAndPort(envAddressValue)
		if err != nil {
			log.Fatal(err)
		}

		addr.host = host
		addr.port = port
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		logLevel = envLogLevelValue
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
