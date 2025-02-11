package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
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
	envAddress         = "ADDRESS"
	envDatabaseDSN     = "DATABASE_DSN"
	envLogLevel        = "LOG_LEVEL"
	envStoreInterval   = "STORE_INTERVAL"
	envFileStoragePath = "FILE_STORAGE_PATH"
	envRestore         = "RESTORE"
	envSecretKey       = "KEY"
	envCryptoKey       = "CRYPTO_KEY"
	envConfigPath      = "CONFIG"
	envTrustedSubnet   = "TRUSTED_SUBNET"
)

type Configuration struct {
	LogLevel        string
	Address         *address
	DatabaseDSN     string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	SecretKey       string
	CryptoKeyPath   string
	ConfigPath      string
	TrustedSubnet   *net.IPNet
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
	flag.StringVar(&c.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&c.LogLevel, "l", "", "Log level")

	flag.Func("i", "Store interval in seconds", func(s string) error {
		var err error

		c.StoreInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return fmt.Errorf("parse store interval error: %w", err)
		}

		return nil
	})

	flag.StringVar(&c.FileStoragePath, "f", "", "File storage path")
	flag.BoolVar(&c.Restore, "r", false, "Restore metrics")
	flag.StringVar(&c.SecretKey, "k", "", "Secret key")
	flag.StringVar(&c.CryptoKeyPath, "crypto-key", "", "Path to private key")
	flag.StringVar(&c.ConfigPath, "c", "", "Path to config JSON file")
	flag.StringVar(&c.ConfigPath, "config", "", "Path to config JSON file")

	flag.Func("t", "Trusted subnet: 192.168.1.0/24", func(s string) error {
		var err error

		_, c.TrustedSubnet, err = net.ParseCIDR(s)
		if err != nil {
			return fmt.Errorf("parse CIDR error: %w", err)
		}

		return nil
	})

	flag.Parse()
}

func (c *Configuration) parseEnvs() {
	if value := os.Getenv(envAddress); value != "" {
		var err error

		c.Address.Host, c.Address.Port, err = parseHostAndPort(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envDatabaseDSN); value != "" {
		c.DatabaseDSN = value
	}

	if value := os.Getenv(envLogLevel); value != "" {
		c.LogLevel = value
	}

	if value := os.Getenv(envStoreInterval); value != "" {
		var err error

		c.StoreInterval, err = stringToDurationInSeconds(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if value := os.Getenv(envFileStoragePath); value != "" {
		c.FileStoragePath = value
	}

	if value := os.Getenv(envRestore); value != "" {
		c.Restore = true
	}

	if value := os.Getenv(envSecretKey); value != "" {
		c.SecretKey = value
	}

	if value := os.Getenv(envCryptoKey); value != "" {
		c.CryptoKeyPath = value
	}

	if value := os.Getenv(envConfigPath); value != "" {
		c.ConfigPath = value
	}

	if value := os.Getenv(envTrustedSubnet); value != "" {
		var err error

		_, c.TrustedSubnet, err = net.ParseCIDR(value)
		if err != nil {
			log.Fatal(err)
		}
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

	if c.DatabaseDSN == "" && parsedConfig.DatabaseDSN != nil {
		c.DatabaseDSN = *parsedConfig.DatabaseDSN
	}

	if c.LogLevel == "" && parsedConfig.LogLevel != nil {
		c.LogLevel = *parsedConfig.LogLevel
	}

	if c.StoreInterval <= 0 && parsedConfig.StoreInterval != nil {
		c.StoreInterval = *parsedConfig.StoreInterval
	}

	if c.FileStoragePath == "" && parsedConfig.FileStoragePath != nil {
		c.FileStoragePath = *parsedConfig.FileStoragePath
	}

	if !c.Restore && parsedConfig.Restore != nil {
		c.Restore = *parsedConfig.Restore
	}

	if c.SecretKey == "" && parsedConfig.SecretKey != nil {
		c.SecretKey = *parsedConfig.SecretKey
	}

	if c.CryptoKeyPath == "" && parsedConfig.CryptoKey != nil {
		c.CryptoKeyPath = *parsedConfig.CryptoKey
	}

	if c.TrustedSubnet == nil && parsedConfig.TrustedSubnet != nil {
		c.TrustedSubnet = parsedConfig.TrustedSubnet
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
