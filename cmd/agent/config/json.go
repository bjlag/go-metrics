package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type jsonConfig struct {
	Address        *address       `json:"address,omitempty"`
	ReportInterval *time.Duration `json:"report_interval,omitempty"`
	PollInterval   *time.Duration `json:"poll_interval,omitempty"`
	CryptoKey      *string        `json:"crypto_key,omitempty"`
	LogLevel       *string        `json:"log_level,omitempty"`
	SecretKey      *string        `json:"key,omitempty"`
	RateLimit      *int           `json:"rate_limit,omitempty"`
}

func (c *jsonConfig) UnmarshalJSON(b []byte) error {
	type alias jsonConfig

	aliasValue := &struct {
		*alias
		Address        *string `json:"address,omitempty"`
		ReportInterval *string `json:"report_interval,omitempty"`
		PollInterval   *string `json:"poll_interval,omitempty"`
	}{
		alias: (*alias)(c),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return fmt.Errorf("unmarshal JSON config error: %w", err)
	}

	if aliasValue.Address != nil {
		host, port, err := parseHostAndPort(*aliasValue.Address)
		if err != nil {
			return fmt.Errorf("parse address error: %w", err)
		}

		c.Address = &address{host, port}
	}

	if aliasValue.ReportInterval != nil {
		interval, err := time.ParseDuration(*aliasValue.ReportInterval)
		if err != nil {
			return fmt.Errorf("parse report_interval error: %w", err)
		}

		c.ReportInterval = &interval
	}

	if aliasValue.PollInterval != nil {
		interval, err := time.ParseDuration(*aliasValue.PollInterval)
		if err != nil {
			return fmt.Errorf("parse poll_interval error: %w", err)
		}

		c.PollInterval = &interval
	}

	return nil
}
