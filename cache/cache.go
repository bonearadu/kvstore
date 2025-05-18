package cache

type Cache interface {
	// Read returns the latest value for a given key from the Cache,
	// and a boolean representing either a cache hit (true) or miss (false).
	Read(key string) (string, bool)

	// Write writes a value to the cache.
	Write(key string, value string)

	// Deletes a key from the cache.
	Delete(key string)
}
