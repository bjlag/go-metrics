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
	defaultPoolInterval   = 2
	defaultReportInterval = 10
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
	pollInterval   = defaultPoolInterval * time.Second
	reportInterval = defaultReportInterval * time.Second
)

func parseFlags() {
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Server address: host:port")
	flag.Func("p", "Poll interval in seconds", func(s string) error {
		var err error

		pollInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return err
		}

		return nil
	})
	flag.Func("r", "Report interval in seconds", func(s string) error {
		var err error

		reportInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return err
		}

		return nil
	})

	flag.Parse()
}

func stringToDurationInSeconds(s string) (time.Duration, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(val) * time.Second, nil
}
