package cache

import (
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	capacity := 10

	cache := NewLRUCache(capacity)

	if cache.capacity != 10 {
		t.Fatalf("Expected capacity %d, got %d.", capacity, cache.capacity)
	}
}

func TestRead(t *testing.T) {
	cache := NewLRUCache(10)

	t.Run("returns correct value for cache hit", func(t *testing.T) {
		key := "abc"
		expectedVal := "def"
		cache.elements[key] = cache.store.PushFront(&entry{key, expectedVal})

		val, ok := cache.Read(key)

		if !ok {
			t.Fatalf("Expected ok = true, got false")
		}
		if val != expectedVal {
			t.Fatalf("Expected value %s, got %s", expectedVal, val)
		}
	})

	t.Run("returns correct value for cache miss", func(t *testing.T) {
		val, ok := cache.Read("no_such_key")

		if val != "" {
			t.Fatalf("Expected empty string for val, got %s", val)
		}
		if ok {
			t.Fatalf("Expected ok = false, got true")
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("updates map state correctly", func(t *testing.T) {
		cache := NewLRUCache(10)

		cache.Write("123", "456")
		cache.Write("some_key", "someVal")
		cache.Write("some_key", "some other val")

		if len(cache.elements) != 2 {
			t.Fatalf("Expected cache store to have length 2, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 2 {
			t.Fatalf("Expected cache queue to have length 2, has actual length %d", cache.store.Len())
		}

		pairs := []struct {
			key string
			val string
		}{
			{"123", "456"},
			{"some_key", "some other val"},
		}

		for _, p := range pairs {
			actual := cache.elements[p.key]
			if actual.Value.(*entry).val != p.val {
				t.Fatalf("Wrong value for key %s. Expected %s, got %s", p.key, p.val, actual.Value.(*entry).val)
			}
		}
	})

	t.Run("triggers eviction when cache is full", func(t *testing.T) {
		cache := NewLRUCache(1)

		cache.Write("1", "11")

		// Cache should be at capacity
		if len(cache.elements) != 1 {
			t.Fatalf("Expected cache store to have length 1, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 1 {
			t.Fatalf("Expected cache queue to have length 1, has actual length %d", cache.store.Len())
		}

		cache.Write("2", "22")

		// Eviction should have been triggered
		if len(cache.elements) != 1 {
			t.Fatalf("Expected cache store to have length 1, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 1 {
			t.Fatalf("Expected cache queue to have length 1, has actual length %d", cache.store.Len())
		}

		v1, ok1 := cache.Read("1")
		v2, ok2 := cache.Read("2")

		if v1 != "" || ok1 != false {
			t.Fatalf("Key \"1\" should have been evicted but is still in the cache")
		}
		if v2 == "" || ok2 == false {
			t.Fatalf("Key \"2\" should not have been evicted")
		}
	})

	t.Run("does not trigger eviction when cache is not full", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Write("1", "11")

		// Cache should be at capacity
		if len(cache.elements) != 1 {
			t.Fatalf("Expected cache store to have length 1, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 1 {
			t.Fatalf("Expected cache queue to have length 1, has actual length %d", cache.store.Len())
		}

		cache.Write("2", "22")

		// Eviction should not have been triggered
		if len(cache.elements) != 2 {
			t.Fatalf("Expected cache store to have length 2, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 2 {
			t.Fatalf("Expected cache queue to have length 2, has actual length %d", cache.store.Len())
		}

		v1, ok1 := cache.Read("1")
		v2, ok2 := cache.Read("2")

		if v1 == "" || ok1 == false {
			t.Fatalf("Key \"1\" should not have been evicted")
		}
		if v2 == "" || ok2 == false {
			t.Fatalf("Key \"2\" should not have been evicted")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes existing key from cache", func(t *testing.T) {
		cache := NewLRUCache(2)
		cache.Write("key1", "val1")
		cache.Write("key2", "val2")

		cache.Delete("key1")

		val, ok := cache.Read("key1")
		if ok {
			t.Fatalf("Expected key1 to be deleted, but it still exists with value %s", val)
		}

		if len(cache.elements) != 1 {
			t.Fatalf("Expected cache store to have length 1 after deletion, has actual length %d", len(cache.elements))
		}
		if cache.store.Len() != 1 {
			t.Fatalf("Expected cache queue to have length 1 after deletion, has actual length %d", cache.store.Len())
		}
	})

	t.Run("does not affect nonexistent key", func(t *testing.T) {
		cache := NewLRUCache(2)
		cache.Write("key1", "val1")

		cache.Delete("nonexistent_key")

		if len(cache.elements) != 1 {
			t.Fatalf("Deleting nonexistent key should not change cache size. Expected 1, got %d", len(cache.elements))
		}
	})
}
