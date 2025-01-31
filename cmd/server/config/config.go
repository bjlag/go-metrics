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
	envAddress         = "ADDRESS"
	envDatabaseDSN     = "DATABASE_DSN"
	envLogLevel        = "LOG_LEVEL"
	envStoreInterval   = "STORE_INTERVAL"
	envFileStoragePath = "FILE_STORAGE_PATH"
	envRestore         = "RESTORE"
	envSecretKey       = "KEY"
	envCryptoKey       = "CRYPTO_KEY"
	envConfigPath      = "CONFIG"
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

	flag.Parse()
}

func (c *Configuration) parseEnvs() {
	if envAddressValue := os.Getenv(envAddress); envAddressValue != "" {
		var err error

		c.Address.Host, c.Address.Port, err = parseHostAndPort(envAddressValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envDatabaseDSNValue := os.Getenv(envDatabaseDSN); envDatabaseDSNValue != "" {
		c.DatabaseDSN = envDatabaseDSNValue
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		c.LogLevel = envLogLevelValue
	}

	if envStoreIntervalValue := os.Getenv(envStoreInterval); envStoreIntervalValue != "" {
		var err error

		c.StoreInterval, err = stringToDurationInSeconds(envStoreIntervalValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envFileStoragePathValue := os.Getenv(envFileStoragePath); envFileStoragePathValue != "" {
		c.FileStoragePath = envFileStoragePathValue
	}

	if envRestoreValue := os.Getenv(envRestore); envRestoreValue != "" {
		c.Restore = true
	}

	if envSecretKeyValue := os.Getenv(envSecretKey); envSecretKeyValue != "" {
		c.SecretKey = envSecretKeyValue
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
