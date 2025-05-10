package kv_store

import (
	"fmt"
	"sync"
)

type InMemoryStore[K comparable, V any] struct {
	mapStore map[K]V
	mu       sync.RWMutex
}

func NewInMemoryStore[K comparable, V any]() *InMemoryStore[K, V] {
	return &InMemoryStore[K, V]{
		mapStore: make(map[K]V),
		mu:       sync.RWMutex{},
	}
}

func (i *InMemoryStore[K, V]) Put(key K, value V) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.mapStore[key] = value

	return nil
}

func (i *InMemoryStore[K, V]) Get(key K) (V, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	value, ok := i.mapStore[key]

	if !ok {
		return value, fmt.Errorf("key not found")
	}
	return value, nil
}

func (i *InMemoryStore[K, V]) Delete(key K) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.mapStore, key)

	return nil
}

func (i *InMemoryStore[K, V]) Entries() ([]Entry[K, V], error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	entries := make([]Entry[K, V], 0, len(i.mapStore))
	for k, v := range i.mapStore {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}

	return entries, nil
}
