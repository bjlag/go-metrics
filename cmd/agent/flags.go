package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost           = "localhost"
	defaultPort           = 8080
	defaultPoolInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

type netAddress struct {
	host string
	port int
}

func (o *netAddress) String() string {
	return fmt.Sprintf("%s:%d", o.host, o.port)
}

func (o *netAddress) Set(value string) error {
	values := strings.Split(value, ":")
	if len(values) != 2 {
		return fmt.Errorf("invalid format")
	}

	port, err := strconv.Atoi(values[1])
	if err != nil {
		return err
	}

	o.host = values[0]
	o.port = port

	return nil
}

var (
	addr = &netAddress{
		host: defaultHost,
		port: defaultPort,
	}
	pollInterval   time.Duration
	reportInterval time.Duration
)

func parseFlags() {
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Server address: host:port")
	flag.DurationVar(&pollInterval, "p", defaultPoolInterval, "Poll interval in seconds")
	flag.DurationVar(&reportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.Parse()
}
