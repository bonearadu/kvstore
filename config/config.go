package config

import "flag"

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port int
}

// ParseFlags parses command-line flags and returns a ServerConfig
func ParseFlags() *ServerConfig {
	config := &ServerConfig{}
	flag.IntVar(&config.Port, "port", 8080, "Port to listen on")
	flag.Parse()
	return config
}
