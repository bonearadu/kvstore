package kv_store

import (
	"sync"

	"github.com/bonearadu/kvstore/cache"
)

type PersistentCachedStore struct {
	store *PersistentStore
	cache cache.Cache
	mu    sync.RWMutex
}

func NewPersistentCachedStore(storeRootPath string, cacheCapacity int) *PersistentCachedStore {
	return &PersistentCachedStore{
		store: NewPersistentStore(storeRootPath),
		cache: cache.NewLRUCache(cacheCapacity),
		mu:    sync.RWMutex{},
	}
}

func (p *PersistentCachedStore) Put(key string, value string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.store.Put(key, value)
	if err != nil {
		return err
	}

	p.cache.Write(key, value)
	return nil
}

func (p *PersistentCachedStore) Get(key string) (string, error) {
	p.mu.RLock()

	val, ok := p.cache.Read(key)
	if ok {
		p.mu.RUnlock()
		return val, nil
	}

	val, err := p.store.Get(key)
	p.mu.RUnlock()

	if err == nil {
		p.mu.Lock()
		defer p.mu.Unlock()

		p.cache.Write(key, val)
		return val, nil
	}

	return val, err
}

func (p *PersistentCachedStore) Delete(key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cache.Delete(key)
	return p.store.Delete(key)
}

func (p *PersistentCachedStore) Entries() ([]Entry, error) {
	return p.store.Entries()
}
