package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	defaultHost           = "localhost"
	defaultPort           = 8080
	defaultPoolInterval   = 2
	defaultReportInterval = 10
	defaultLogLevel       = "info"
	defaultSecretKey      = "secretkey"

	envAddressKey        = "ADDRESS"
	envPollIntervalKey   = "POLL_INTERVAL"
	envReportIntervalKey = "REPORT_INTERVAL"
	envLogLevel          = "LOG_LEVEL"
	envSecretKey         = "KEY"
)

var (
	addr = &netAddress{
		host: defaultHost,
		port: defaultPort,
	}

	pollInterval   = defaultPoolInterval * time.Second
	reportInterval = defaultReportInterval * time.Second

	logLevel  string
	secretKey string
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
	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.StringVar(&secretKey, "k", defaultSecretKey, "Secret key")

	flag.Parse()
}

func parseEnvs() {
	var (
		host string
		port int
		err  error
	)

	if address := os.Getenv(envAddressKey); address != "" {
		host, port, err = parseHostAndPort(address)
		if err != nil {
			log.Fatal(err)
		}

		addr.host = host
		addr.port = port
	}

	if envPollInterval := os.Getenv(envPollIntervalKey); envPollInterval != "" {
		pollInterval, err = stringToDurationInSeconds(envPollInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envReportInterval := os.Getenv(envReportIntervalKey); envReportInterval != "" {
		reportInterval, err = stringToDurationInSeconds(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		logLevel = envLogLevelValue
	}

	if envSecretKeyValue := os.Getenv(envSecretKey); envSecretKeyValue != "" {
		secretKey = envSecretKeyValue
	}
}

func stringToDurationInSeconds(s string) (time.Duration, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(val) * time.Second, nil
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
