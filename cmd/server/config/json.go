package config

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type jsonConfig struct {
	AddressHTTP     *address       `json:"address,omitempty"`
	AddressRPC      *address       `json:"address_rpc,omitempty"`
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
		AddressHTTP   *string `json:"address,omitempty"`
		AddressRPC    *string `json:"address_rpc,omitempty"`
		StoreInterval *string `json:"store_interval,omitempty"`
		TrustedSubnet *string `json:"trusted_subnet,omitempty"`
	}{
		alias: (*alias)(c),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return fmt.Errorf("unmarshal JSON config error: %w", err)
	}

	if aliasValue.AddressHTTP != nil && *aliasValue.AddressHTTP != "" {
		host, port, err := parseHostAndPort(*aliasValue.AddressHTTP)
		if err != nil {
			return fmt.Errorf("parse address error: %w", err)
		}

		c.AddressHTTP = &address{host, port}
	}

	if aliasValue.AddressRPC != nil && *aliasValue.AddressRPC != "" {
		host, port, err := parseHostAndPort(*aliasValue.AddressRPC)
		if err != nil {
			return fmt.Errorf("parse address error: %w", err)
		}

		c.AddressRPC = &address{host, port}
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
