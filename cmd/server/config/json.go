package config

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type jsonConfig struct {
	Address         *address       `json:"address,omitempty"`
	Restore         *bool          `json:"restore,omitempty"`
	StoreInterval   *time.Duration `json:"store_interval,omitempty"`
	StoreFile       *string        `json:"store_file,omitempty"`
	DatabaseDSN     *string        `json:"database_dsn,omitempty"`
	CryptoKey       *string        `json:"crypto_key,omitempty"`
	LogLevel        *string        `json:"log_level,omitempty"`
	FileStoragePath *string        `json:"file_storage_path,omitempty"`
	SecretKey       *string        `json:"key,omitempty"`
	TrustedSubnet   *net.IPNet     `json:"trusted_subnet,omitempty"`
}

func (c *jsonConfig) UnmarshalJSON(b []byte) error {
	type alias jsonConfig

	aliasValue := &struct {
		*alias
		Address       *string `json:"address,omitempty"`
		StoreInterval *string `json:"store_interval,omitempty"`
		TrustedSubnet *string `json:"trusted_subnet,omitempty"`
	}{
		alias: (*alias)(c),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return fmt.Errorf("unmarshal JSON config error: %w", err)
	}

	if aliasValue.Address != nil && *aliasValue.Address != "" {
		host, port, err := parseHostAndPort(*aliasValue.Address)
		if err != nil {
			return fmt.Errorf("parse address error: %w", err)
		}

		c.Address = &address{host, port}
	}

	if aliasValue.StoreInterval != nil && *aliasValue.StoreInterval != "" {
		interval, err := time.ParseDuration(*aliasValue.StoreInterval)
		if err != nil {
			return fmt.Errorf("parse store_interval error: %w", err)
		}

		c.StoreInterval = &interval
	}

	if aliasValue.TrustedSubnet != nil && *aliasValue.TrustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(*aliasValue.TrustedSubnet)
		if err != nil {
			return fmt.Errorf("parse CIDR error: %w", err)
		}

		c.TrustedSubnet = ipNet
	}

	return nil
}
