package config

import "flag"

type StoreImpl int

const (
	InMemory StoreImpl = iota
	InMemorySnap
	Persistent
	PersistentCached
)

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port          int
	Mode          StoreImpl
	StorePath     string
	CacheCapacity int
}

// ParseFlags parses command-line flags and returns a ServerConfig
func ParseFlags() *ServerConfig {
	config := &ServerConfig{}
	var mode int

	flag.IntVar(&config.Port, "port", 8080, "Port to listen on")
	flag.IntVar(&mode, "mode", 0,
		"The key-value store implementation to use. 0 = In-Memory map, 1 = In-Memory with snapshotting"+
			"2 = Persistent KV store, 3 = Persistent KV store with caching")
	flag.StringVar(&config.StorePath, "store_path", "", "The path for the persistent storage or snapshots, if used")
	flag.IntVar(&config.CacheCapacity, "cache_capacity", 100,
		"The size of the cache for the persistent cached storage, if used")
	flag.Parse()
	config.Mode = StoreImpl(mode)

	return config
}
