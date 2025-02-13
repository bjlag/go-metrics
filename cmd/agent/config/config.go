package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type address struct {
	Host string
	Port int
}

func (o *address) String() string {
	return fmt.Sprintf("%s:%d", o.Host, o.Port)
}

func (o *address) Set(value string) error {
	var err error

	o.Host, o.Port, err = parseHostAndPort(value)
	if err != nil {
		return err
	}

	return nil
}

const (
	envAddressHTTP    = "ADDRESS"
	envAddressRPC     = "ADDRESS_RPC"
	envPollInterval   = "POLL_INTERVAL"
	envReportInterval = "REPORT_INTERVAL"
	envLogLevel       = "LOG_LEVEL"
	envSecret         = "KEY"
	envRateLimit      = "RATE_LIMIT"
	envCrypto         = "CRYPTO_KEY"
	envConfigPath     = "CONFIG"
)

type Configuration struct {
	LogLevel       string
	AddressHTTP    *address
	AddressRPC     *address
	ReportInterval time.Duration
	PollInterval   time.Duration
	CryptoKeyPath  string
	SecretKey      string
	RateLimit      int
	ConfigPath     string
}

func LoadConfig() *Configuration {
	c := &Configuration{
		AddressHTTP: &address{},
		AddressRPC:  &address{},
	}

	c.parseFlags()
	c.parseEnvs()
	c.parseJSONConfig()

	return c
}

func (c *Configuration) parseFlags() {
	_ = flag.Value(c.AddressHTTP)
	_ = flag.Value(c.AddressRPC)

	flag.Var(c.AddressHTTP, "a", "Server HTTP address: host:port")
	flag.Var(c.AddressRPC, "address-rpc", "Server RPC address: host:port")
	flag.Func("p", "Poll interval in seconds", func(s string) error {
		var err error

		c.PollInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return fmt.Errorf("parse poll interval error: %w", err)
		}

		return nil
	})
	flag.Func("r", "Report interval in seconds", func(s string) error {
		var err error

		c.ReportInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return fmt.Errorf("parse report interval error: %w", err)
		}

		return nil
	})
	flag.StringVar(&c.LogLevel, "log", "", "Log level")
	flag.StringVar(&c.SecretKey, "k", "", "Secret key")
	flag.IntVar(&c.RateLimit, "l", 0, "Rate limit")
	flag.StringVar(&c.CryptoKeyPath, "crypto-key", "", "Path to public key")
	flag.StringVar(&c.ConfigPath, "c", "", "Path to config JSON file")
	flag.StringVar(&c.ConfigPath, "config", "", "Path to config JSON file")

	flag.Parse()
}

func (c *Configuration) parseEnvs() {
	var err error

	if value := os.Getenv(envAddressHTTP); value != "" {
		c.AddressHTTP.Host, c.AddressHTTP.Port, err = parseHostAndPort(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envAddressRPC); value != "" {
		c.AddressRPC.Host, c.AddressRPC.Port, err = parseHostAndPort(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envPollInterval); value != "" {
		c.PollInterval, err = stringToDurationInSeconds(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envReportInterval); value != "" {
		c.ReportInterval, err = stringToDurationInSeconds(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envLogLevel); value != "" {
		c.LogLevel = value
	}

	if value := os.Getenv(envSecret); value != "" {
		c.SecretKey = value
	}

	if value := os.Getenv(envRateLimit); value != "" {
		c.RateLimit, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envCrypto); value != "" {
		c.CryptoKeyPath = value
	}

	if value := os.Getenv(envConfigPath); value != "" {
		c.ConfigPath = value
	}
}

func (c *Configuration) parseJSONConfig() {
	if c.ConfigPath == "" {
		return
	}

	configFile, err := os.Open(c.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = configFile.Close()
	}()

	var parsedConfig jsonConfig
	err = json.NewDecoder(configFile).Decode(&parsedConfig)
	if err != nil {
		log.Fatal(err)
	}

	if c.AddressHTTP.Host == "" && c.AddressHTTP.Port <= 0 && parsedConfig.AddressHTTP != nil {
		c.AddressHTTP = parsedConfig.AddressHTTP
	}

	if c.AddressRPC.Host == "" && c.AddressRPC.Port <= 0 && parsedConfig.AddressRPC != nil {
		c.AddressRPC = parsedConfig.AddressRPC
	}

	if c.LogLevel == "" && parsedConfig.LogLevel != nil {
		c.LogLevel = *parsedConfig.LogLevel
	}

	if c.PollInterval <= 0 && parsedConfig.PollInterval != nil {
		c.PollInterval = *parsedConfig.PollInterval
	}

	if c.ReportInterval <= 0 && parsedConfig.ReportInterval != nil {
		c.ReportInterval = *parsedConfig.ReportInterval
	}

	if c.SecretKey == "" && parsedConfig.SecretKey != nil {
		c.SecretKey = *parsedConfig.SecretKey
	}

	if c.CryptoKeyPath == "" && parsedConfig.CryptoKey != nil {
		c.CryptoKeyPath = *parsedConfig.CryptoKey
	}

	if c.RateLimit <= 0 && parsedConfig.RateLimit != nil {
		c.RateLimit = *parsedConfig.RateLimit
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
