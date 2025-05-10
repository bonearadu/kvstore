package kv_store

// KeyValueStore defines the interface for a generic key-value store.
// K represents the key type which must be comparable.
// V represents the value type which can be any type.
type KeyValueStore[K comparable, V any] interface {
	// Put stores the given value associated with the given key.
	// If the key already exists, its value is updated.
	// Returns an error if the operation fails.
	Put(key K, value V) error

	// Get retrieves the value associated with the given key.
	// Returns the value and nil error if the key exists.
	// Returns a zero value and an error if the key doesn't exist or if the operation fails.
	Get(key K) (V, error)

	// Delete removes the key-value pair for the given key.
	// Returns nil if the key was successfully deleted or didn't exist.
	// Returns an error if the operation fails.
	Delete(key K) error

	// Entries returns all key-value pairs in the store.
	// Returns a slice of Entry structs and nil error on success.
	// Returns an empty slice and an error if the operation fails.
	Entries() ([]Entry[K, V], error)
}

// Entry represents a key-value pair in the store.
type Entry[K comparable, V any] struct {
	// Key is the identifier for the value.
	Key K
	// Value is the data associated with the key.
	Value V
}
