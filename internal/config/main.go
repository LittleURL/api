package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

type Config struct {
	Auth0  auth0Config
	Queues queues
	Tables tables
}

type auth0Config struct {
	Domain       string
	ClientID     string
	ClientSecret string
}

type queues struct{}

type tables struct {
	Domains     string
	DomainUsers string
}

func Load() *Config {
	config := &Config{}

	// parse config
	err := configor.New(&configor.Config{
		ErrorOnUnmatchedKeys: true,
	}).Load(config)

	// log error and kill server if config is invalid
	if err != nil {
		fmt.Printf("Failed to load config: %s", err)
		os.Exit(1)
	}

	return config
}
