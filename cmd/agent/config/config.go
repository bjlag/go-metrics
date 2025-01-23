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
	envAddressKey        = "ADDRESS"
	envPollIntervalKey   = "POLL_INTERVAL"
	envReportIntervalKey = "REPORT_INTERVAL"
	envLogLevel          = "LOG_LEVEL"
	envSecretKey         = "KEY"
	envRateLimitKey      = "RATE_LIMIT"
	envCryptoKey         = "CRYPTO_KEY"
	envConfigPath        = "CONFIG"
)

type Configuration struct {
	LogLevel       string
	Address        *address
	ReportInterval time.Duration
	PollInterval   time.Duration
	CryptoKeyPath  string
	SecretKey      string
	RateLimit      int
	ConfigPath     string
}

func LoadConfig() *Configuration {
	c := &Configuration{
		Address: &address{},
	}

	c.parseFlags()
	c.parseEnvs()
	c.parseJSONConfig()

	return c
}

func (c *Configuration) parseFlags() {
	_ = flag.Value(c.Address)

	flag.Var(c.Address, "a", "Server address: host:port")
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

	if envAddressValue := os.Getenv(envAddressKey); envAddressValue != "" {
		c.Address.Host, c.Address.Port, err = parseHostAndPort(envAddressValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envPollInterval := os.Getenv(envPollIntervalKey); envPollInterval != "" {
		c.PollInterval, err = stringToDurationInSeconds(envPollInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envReportInterval := os.Getenv(envReportIntervalKey); envReportInterval != "" {
		c.ReportInterval, err = stringToDurationInSeconds(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		c.LogLevel = envLogLevelValue
	}

	if envSecretKeyValue := os.Getenv(envSecretKey); envSecretKeyValue != "" {
		c.SecretKey = envSecretKeyValue
	}

	if envRateLimitValue := os.Getenv(envRateLimitKey); envRateLimitValue != "" {
		c.RateLimit, err = strconv.Atoi(envRateLimitValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envCryptoKeyValue := os.Getenv(envCryptoKey); envCryptoKeyValue != "" {
		c.CryptoKeyPath = envCryptoKeyValue
	}

	if envConfigPathValue := os.Getenv(envConfigPath); envConfigPathValue != "" {
		c.ConfigPath = envConfigPathValue
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

	if c.Address.Host == "" && c.Address.Port <= 0 && parsedConfig.Address != nil {
		c.Address = parsedConfig.Address
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
