package kv_store

import (
	"os"
	"testing"
)

func TestNewPersistentStore(t *testing.T) {
	t.Run("normal creation", func(t *testing.T) {
		storeRoot, err := os.MkdirTemp("", "persistent_store_test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(storeRoot)

		store := NewPersistentStore(storeRoot)
		if store == nil {
			t.Fatal("Expected NewPersistentStore to return a non-nil store")
		}
		if store.storeRoot != storeRoot {
			t.Fatalf("Expected storeRoot to be %s, got %s", storeRoot, store.storeRoot)
		}
	})

	t.Run("non-existent root path", func(t *testing.T) {
		const tempPath = "non_existent_path"
		defer os.RemoveAll(tempPath)

		store := NewPersistentStore(tempPath)

		if store == nil {
			t.Fatal("Expected NewPersistentStore to return a non-nil store")
		}
		if store.storeRoot != tempPath {
			t.Fatalf("Expected storeRoot to be %s, got %s", tempPath, store.storeRoot)
		}
	})
}

func TestPut(t *testing.T) {
	storeRoot, err := os.MkdirTemp("", "persistent_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(storeRoot)

	store := NewPersistentStore(storeRoot)

	t.Run("new key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != "value1" {
			t.Fatalf("Expected value 'value1', got '%v'", value)
		}
	})

	t.Run("update existing key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		err = store.Put("key1", "value2")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != "value2" {
			t.Fatalf("Expected value 'value2', got '%v'", value)
		}
	})
}

func TestGet(t *testing.T) {
	storeRoot, err := os.MkdirTemp("", "persistent_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(storeRoot)

	store := NewPersistentStore(storeRoot)

	t.Run("existing key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		value, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if value != "value1" {
			t.Fatalf("Expected value 'value1', got '%v'", value)
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		_, err := store.Get("nonexistent")
		if err == nil {
			t.Fatal("Expected error for non-existent key, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	storeRoot, err := os.MkdirTemp("", "persistent_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(storeRoot)

	store := NewPersistentStore(storeRoot)

	t.Run("existing key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		err = store.Delete("key1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = store.Get("key1")
		if err == nil {
			t.Fatal("Expected error after deletion, got nil")
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		err := store.Delete("nonexistent")
		if err != nil {
			t.Fatalf("Expected no error when deleting non-existent key, got %v", err)
		}
	})
}

func TestPersistentStoreEntries(t *testing.T) {
	storeRoot, err := os.MkdirTemp("", "persistent_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(storeRoot)

	store := NewPersistentStore(storeRoot)

	err = store.Put("key1", "value1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	err = store.Put("key2", "value2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	entries, err := store.Entries()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	entryMap := make(map[string]string)
	for _, entry := range entries {
		entryMap[entry.Key] = entry.Value
	}

	if entryMap["key1"] != "value1" || entryMap["key2"] != "value2" {
		t.Fatalf("Expected entries {key1: value1, key2: value2}, got %+v", entryMap)
	}
}

func TestPersistentStorePersistence(t *testing.T) {
	storeRoot, err := os.MkdirTemp("", "persistent_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(storeRoot)

	store1 := NewPersistentStore(storeRoot)
	err = store1.Put("key1", "value1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	store2 := NewPersistentStore(storeRoot)
	value, err := store2.Get("key1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "value1" {
		t.Fatalf("Expected value 'value1', got '%v'", value)
	}
}
