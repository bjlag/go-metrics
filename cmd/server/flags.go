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
	defaultHost            = "localhost"
	defaultPort            = 8080
	defaultLogLevel        = "info"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "data/metrics.json"
	defaultRestore         = true

	envAddress         = "ADDRESS"
	envDatabaseDSN     = "DATABASE_DSN"
	envLogLevel        = "LOG_LEVEL"
	envStoreInterval   = "STORE_INTERVAL"
	envFileStoragePath = "FILE_STORAGE_PATH"
	envRestore         = "RESTORE"
	envSecretKey       = "KEY"
	envCryptoKey       = "CRYPTO_KEY"
)

var (
	addr = &netAddress{
		host: defaultHost,
		port: defaultPort,
	}

	databaseDSN     string
	logLevel        string
	storeInterval   = defaultStoreInterval * time.Second
	fileStoragePath string
	restore         bool
	secretKey       string
	cryptoKeyPath   string
)

func parseFlags() {
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Server address: host:port")
	flag.StringVar(&databaseDSN, "d", "", "Database DSN")
	flag.StringVar(&logLevel, "l", defaultLogLevel, "Log level")
	flag.Func("i", "Store interval in seconds", func(s string) error {
		var err error

		storeInterval, err = stringToDurationInSeconds(s)
		if err != nil {
			return err
		}

		return nil
	})
	flag.StringVar(&fileStoragePath, "f", defaultFileStoragePath, "File storage path")
	flag.BoolVar(&restore, "r", defaultRestore, "Restore metrics")
	flag.StringVar(&secretKey, "k", "", "Secret key")
	flag.StringVar(&cryptoKeyPath, "crypto-key", "", "Path to private key")

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

	if envDatabaseDSNValue := os.Getenv(envDatabaseDSN); envDatabaseDSNValue != "" {
		databaseDSN = envDatabaseDSNValue
	}

	if envLogLevelValue := os.Getenv(envLogLevel); envLogLevelValue != "" {
		logLevel = envLogLevelValue
	}

	if envStoreIntervalValue := os.Getenv(envStoreInterval); envStoreIntervalValue != "" {
		var err error

		storeInterval, err = stringToDurationInSeconds(envStoreIntervalValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	if envFileStoragePathValue := os.Getenv(envFileStoragePath); envFileStoragePathValue != "" {
		fileStoragePath = envFileStoragePathValue
	}

	if envRestoreValue := os.Getenv(envRestore); envRestoreValue != "" {
		restore = true
	}

	if envSecretKeyValue := os.Getenv(envSecretKey); envSecretKeyValue != "" {
		secretKey = envSecretKeyValue
	}

	if envCryptoKeyValue := os.Getenv(envCryptoKey); envCryptoKeyValue != "" {
		cryptoKeyPath = envCryptoKeyValue
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

func stringToDurationInSeconds(s string) (time.Duration, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(val) * time.Second, nil
}
