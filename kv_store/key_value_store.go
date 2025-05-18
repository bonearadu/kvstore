package kv_store

type KeyValueStore interface {
	// Put stores the given value associated with the given key.
	// If the key already exists, its value is updated.
	// Returns an error if the operation fails.
	Put(key string, value string) error

	// Get retrieves the value associated with the given key.
	// Returns the value and nil error if the key exists.
	// Returns a zero value and an error if the key doesn't exist or if the operation fails.
	Get(key string) (string, error)

	// Delete removes the key-value pair for the given key.
	// Returns nil if the key was successfully deleted or didn't exist.
	// Returns an error if the operation fails.
	Delete(key string) error

	// Entries returns all key-value pairs in the store.
	// Returns a slice of Entry structs and nil error on success.
	// Returns an empty slice and an error if the operation fails.
	Entries() ([]Entry, error)
}

// Entry represents a key-value pair in the store.
type Entry struct {
	// Key is the identifier for the value.
	Key string
	// Value is the data associated with the key.
	Value string
}
