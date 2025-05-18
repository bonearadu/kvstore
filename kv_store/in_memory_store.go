package kv_store

import (
	"fmt"
	"sync"
)

type InMemoryStore struct {
	mapStore map[string]string
	mu       sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		mapStore: make(map[string]string),
		mu:       sync.RWMutex{},
	}
}

func (i *InMemoryStore) Put(key string, value string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.mapStore[key] = value

	return nil
}

func (i *InMemoryStore) Get(key string) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	value, ok := i.mapStore[key]

	if !ok {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

func (i *InMemoryStore) Delete(key string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.mapStore, key)

	return nil
}

func (i *InMemoryStore) Entries() ([]Entry, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	entries := make([]Entry, 0, len(i.mapStore))
	for k, v := range i.mapStore {
		entries = append(entries, Entry{Key: k, Value: v})
	}

	return entries, nil
}
