package kv_store

import (
	"os"
	"testing"
	"time"
)

func TestNewInMemorySnapshottedStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 10)
	if store == nil {
		t.Fatal("Expected NewInMemorySnapshottedStore to return a non-nil store")
	}
	if store.snapshotsRoot != tempDir {
		t.Fatalf("Expected snapshotsRoot to be %s, got %s", tempDir, store.snapshotsRoot)
	}
	if store.snapshotFreqSeconds != 10 {
		t.Fatalf("Expected snapshotFreqSeconds to be 10, got %d", store.snapshotFreqSeconds)
	}
}

func TestInMemorySnapshottedStorePut(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 10)

	t.Run("basic put", func(t *testing.T) {
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

	t.Run("update existing key", func(t *testing.T) {
		err := store.Put("key1", "value2")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		val, err := store.Get("key1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "value2" {
			t.Fatalf("Expected value2, got %s", val)
		}
	})
}

func TestInMemorySnapshottedStoreGet(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 10)

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
			t.Fatalf("Expected value1, got %s", val)
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		_, err := store.Get("nonexistent")
		if err == nil {
			t.Fatal("Expected error for non-existent key, got nil")
		}
	})
}

func TestInMemorySnapshottedStoreDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 10)

	t.Run("existing key", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		err = store.Delete("key1")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
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

func TestInMemorySnapshottedStoreEntries(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 10)

	t.Run("get all entries", func(t *testing.T) {
		// Add entries
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
		err = store.Put("key2", "value2")
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

		expectedKeys := []string{"key1", "key2"}
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

func TestInMemorySnapshottedStoreSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "snapshotted_store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := NewInMemorySnapshottedStore(tempDir, 1)

	t.Run("save snapshot", func(t *testing.T) {
		err := store.Put("key1", "value1")
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}

		time.Sleep(2 * time.Second) // Wait for the snapshot to be saved

		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read snapshots directory: %v", err)
		}
		if len(files) != 1 {
			t.Fatalf("Expected 1 snapshot file, got %d", len(files))
		}
	})

	t.Run("load snapshot", func(t *testing.T) {
		store2 := NewInMemorySnapshottedStore(tempDir, 1)

		val, err := store2.Get("key1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if val != "value1" {
			t.Fatalf("Expected value1, got %s", val)
		}
	})
}
