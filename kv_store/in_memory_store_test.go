package kv_store

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewInMemoryStore(t *testing.T) {
	t.Run("returns non-nil store", func(t *testing.T) {
		store := NewInMemoryStore()
		if store == nil {
			t.Fatal("Expected NewInMemoryStore to return a non-nil store")
		}
	})

	t.Run("initializes map store", func(t *testing.T) {
		store := NewInMemoryStore()
		if store.mapStore == nil {
			t.Fatal("Expected mapStore to be initialized")
		}
	})
}

func TestEntries(t *testing.T) {
	t.Run("empty store", func(t *testing.T) {
		store := NewInMemoryStore()

		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("Expected empty entries, got %d entries", len(entries))
		}
	})

	t.Run("populated store", func(t *testing.T) {
		store := NewInMemoryStore()

		// Setup: Add some key-value pairs
		store.Put("key1", "1")
		store.Put("key2", "2")
		store.Put("key3", "3")

		// Test: Get all entries
		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(entries) != 3 {
			t.Fatalf("Expected 3 entries, got %d entries", len(entries))
		}

		// Create a map for easier verification
		entryMap := make(map[string]string)
		for _, entry := range entries {
			entryMap[entry.Key] = entry.Value
		}

		// Verify all entries are present
		expectedValues := map[string]string{
			"key1": "1",
			"key2": "2",
			"key3": "3",
		}

		for key, expectedValue := range expectedValues {
			value, exists := entryMap[key]
			if !exists {
				t.Fatalf("Expected entry with key %s, but it was not found", key)
			}
			if value != expectedValue {
				t.Fatalf("Expected value %s for key %s, got %s", expectedValue, key, value)
			}
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("concurrent writers and readers", func(t *testing.T) {
		store := NewInMemoryStore()
		const goroutines = 10
		const operationsPerGoroutine = 100

		var wg sync.WaitGroup
		wg.Add(goroutines * 2) // For readers and writers

		// Launch writer goroutines
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					err := store.Put(key, fmt.Sprintf("%d", id*1000+j))
					if err != nil {
						t.Errorf("Error in Put: %v", err)
					}
				}
			}(i)
		}

		// Launch reader goroutines
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					// Try to read keys written by the writer with the same ID
					key := fmt.Sprintf("key-%d-%d", id, j)
					_, _ = store.Get(key) // Errors are expected as writers might not have written yet
				}
			}(i)
		}

		wg.Wait()

		// Verify the final state
		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Error in Entries: %v", err)
		}

		expectedCount := goroutines * operationsPerGoroutine
		if len(entries) != expectedCount {
			t.Fatalf("Expected %d entries, got %d", expectedCount, len(entries))
		}
	})
}

func TestTypes(t *testing.T) {
	t.Run("string keys and string values", func(t *testing.T) {
		store := NewInMemoryStore()

		err := store.Put("hello", "world")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		value, err := store.Get("hello")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != "world" {
			t.Fatalf("Expected value 'world', got '%v'", value)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		store := NewInMemoryStore()

		err := store.Put("", "empty key")
		if err != nil {
			t.Fatalf("Expected no error for empty key, got %v", err)
		}

		value, err := store.Get("")
		if err != nil {
			t.Fatalf("Expected no error for empty key, got %v", err)
		}
		if value != "empty key" {
			t.Fatalf("Expected value 'empty key', got '%v'", value)
		}
	})

	t.Run("empty string value", func(t *testing.T) {
		store := NewInMemoryStore()

		err := store.Put("empty value", "")
		if err != nil {
			t.Fatalf("Expected no error for empty value, got %v", err)
		}

		value, err := store.Get("empty value")
		if err != nil {
			t.Fatalf("Expected no error for key with empty value, got %v", err)
		}
		if value != "" {
			t.Fatalf("Expected empty value, got '%v'", value)
		}
	})
}
