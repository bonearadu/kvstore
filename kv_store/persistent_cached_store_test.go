package kv_store

import (
	"os"
	"testing"
)

func TestPersistentCachedStorePut(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "persistent_cached_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewPersistentCachedStore(tempDir, 2)

	t.Run("basic put and get", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		val, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "value1" {
			t.Fatalf("Expected value1, got %s", val)
		}
	})

	t.Run("cache eviction", func(t *testing.T) {
		// Fill the cache
		err := store.Put("key3", "value3")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
		err = store.Put("key4", "value4")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		// This should evict key3 since capacity is 2
		err = store.Put("key5", "value5")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		// key3 should be evicted from cache but still exist in persistent store
		_, err = store.Get("key3")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
	})
}

func TestPersistentCachedStoreGet(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "persistent_cached_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewPersistentCachedStore(tempDir, 2)

	t.Run("existing key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		val, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected value1, got %s", val)
		}

		// Test cache hit
		val, err = store.Get("key1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected value1, got %s", val)
		}
	})

	t.Run("non-existent key", func(t *testing.T) {
		val, err := store.Get("nonexistent")
		if err == nil {
			t.Fatalf("Expected error for nonexistent key, got nil")
		}
		if val != "" {
			t.Errorf("Expected empty value, got %s", val)
		}
	})
}

func TestPersistentCachedStoreDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "persistent_cached_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewPersistentCachedStore(tempDir, 2)

	t.Run("delete existing key", func(t *testing.T) {
		err := store.Put("key2", "value2")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		err = store.Delete("key2")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		val, err := store.Get("key2")
		if err == nil {
			t.Fatalf("Expected error for deleted key, got nil")
		}
		if val != "" {
			t.Errorf("Expected empty value, got %s", val)
		}
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		err := store.Delete("nonexistent")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})
}

func TestPersistentCachedStoreEntries(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "persistent_cached_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewPersistentCachedStore(tempDir, 2)

	t.Run("get all entries", func(t *testing.T) {
		// Add entries
		err := store.Put("key6", "value6")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
		err = store.Put("key7", "value7")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		entries, err := store.Entries()
		if err != nil {
			t.Fatalf("Entries failed: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(entries))
		}

		expectedKeys := []string{"key6", "key7"}
		for _, key := range expectedKeys {
			found := false
			for _, entry := range entries {
				if entry.Key == key {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected key %s not found in entries", key)
			}
		}
	})
}
