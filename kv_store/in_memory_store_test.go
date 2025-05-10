package kv_store

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewInMemoryStore(t *testing.T) {
	t.Run("returns non-nil store", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		if store == nil {
			t.Fatal("Expected NewInMemoryStore to return a non-nil store")
		}
	})

	t.Run("initializes map store", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		if store.mapStore == nil {
			t.Fatal("Expected mapStore to be initialized")
		}
	})
}

func TestPut(t *testing.T) {
	t.Run("new key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		err := store.Put("key1", 42)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != 42 {
			t.Fatalf("Expected value 42, got %v", value)
		}
	})

	t.Run("update existing key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		// Setup: Add a key first
		err := store.Put("key1", 42)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		
		// Test: Update the key
		err = store.Put("key1", 100)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != 100 {
			t.Fatalf("Expected value 100, got %v", value)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("non-existing key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		_, err := store.Get("nonexistent")
		if err == nil {
			t.Fatal("Expected error for non-existent key, got nil")
		}
	})
	
	t.Run("existing key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		// Setup: Add a key first
		err := store.Put("key1", 42)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		
		// Test: Get the key
		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != 42 {
			t.Fatalf("Expected value 42, got %v", value)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		// Setup: Add a key first
		err := store.Put("key1", 42)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		
		// Test: Delete the key
		err = store.Delete("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Verify: Key should be gone
		_, err = store.Get("key1")
		if err == nil {
			t.Fatal("Expected error after deletion, got nil")
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		// Test: Delete a non-existing key (should be idempotent)
		err := store.Delete("nonexistent")
		if err != nil {
			t.Fatalf("Expected no error when deleting non-existent key, got %v", err)
		}
	})
}

func TestEntries(t *testing.T) {
	t.Run("empty store", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("Expected empty entries, got %d entries", len(entries))
		}
	})

	t.Run("populated store", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
		
		// Setup: Add some key-value pairs
		store.Put("key1", 1)
		store.Put("key2", 2)
		store.Put("key3", 3)
		
		// Test: Get all entries
		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(entries) != 3 {
			t.Fatalf("Expected 3 entries, got %d entries", len(entries))
		}
		
		// Create a map for easier verification
		entryMap := make(map[string]int)
		for _, entry := range entries {
			entryMap[entry.Key] = entry.Value
		}
		
		// Verify all entries are present
		expectedValues := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		}
		
		for key, expectedValue := range expectedValues {
			value, exists := entryMap[key]
			if !exists {
				t.Fatalf("Expected entry with key %s, but it was not found", key)
			}
			if value != expectedValue {
				t.Fatalf("Expected value %d for key %s, got %d", expectedValue, key, value)
			}
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("concurrent writers and readers", func(t *testing.T) {
		store := NewInMemoryStore[string, int]()
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
					err := store.Put(key, id*1000+j)
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
		store := NewInMemoryStore[string, string]()
		
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
	
	t.Run("int keys and struct values", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		
		store := NewInMemoryStore[int, Person]()
		
		err := store.Put(1, Person{Name: "Alice", Age: 30})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		value, err := store.Get(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value.Name != "Alice" || value.Age != 30 {
			t.Fatalf("Expected Person{Name: 'Alice', Age: 30}, got %v", value)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		store := NewInMemoryStore[string, string]()
		
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
		store := NewInMemoryStore[string, string]()
		
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
