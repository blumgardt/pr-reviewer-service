package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HTTP     HTTPConfig     `toml:"http"`
	Postgres PostgresConfig `toml:"postgres"`
}

type HTTPConfig struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type PostgresConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"sslmode"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %s", err)
	}

	return &cfg, nil
}
