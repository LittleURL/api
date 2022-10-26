package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

type Config struct {
	Queues queues
	Tables tables
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
		ENVPrefix:            "-",
	}).Load(config)

	// log error and kill server if config is invalid
	if err != nil {
		fmt.Printf("Failed to load config: %s", err)
		os.Exit(1)
	}

	return config
}
