package kv_store

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

type PersistentStore struct {
	storeRoot   string
	keyMutexMap map[string]*sync.RWMutex
	mapMutex    sync.RWMutex
}

const fileMode = 0777

func NewPersistentStore(storeRootPath string) *PersistentStore {
	if _, err := os.Stat(storeRootPath); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(storeRootPath, fileMode)
	}

	return &PersistentStore{
		storeRoot:   storeRootPath,
		keyMutexMap: make(map[string]*sync.RWMutex),
		mapMutex:    sync.RWMutex{},
	}
}

func (p *PersistentStore) getMutex(key string) *sync.RWMutex {
	p.mapMutex.RLock()
	m, ok := p.keyMutexMap[key]
	p.mapMutex.RUnlock()

	if !ok {
		m = &sync.RWMutex{}

		p.mapMutex.Lock()
		p.keyMutexMap[key] = m
		p.mapMutex.Unlock()
	}

	return m
}

func (p *PersistentStore) removeMutex(key string) {
	p.mapMutex.Lock()
	delete(p.keyMutexMap, key)
	p.mapMutex.Unlock()
}

func (p *PersistentStore) Put(key string, value string) error {
	mu := p.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	bytes := []byte(value)
	err := os.WriteFile(path.Join(p.storeRoot, key), bytes, fileMode)

	if err != nil {
		return fmt.Errorf("error writing value for key %s", key)
	}
	return nil
}

func (p *PersistentStore) Get(key string) (string, error) {
	mu := p.getMutex(key)
	mu.RLock()
	defer mu.RUnlock()

	return p.getUnsafe(key)
}

func (p *PersistentStore) getUnsafe(key string) (string, error) {
	bytes, err := os.ReadFile(path.Join(p.storeRoot, key))

	if err != nil {
		return "", fmt.Errorf("key not found")
	}

	return string(bytes), nil
}

func (p *PersistentStore) Delete(key string) error {
	mu := p.getMutex(key)
	mu.Lock()
	defer func() {
		mu.Unlock()
		p.removeMutex(key)
	}()

	err := os.Remove(path.Join(p.storeRoot, key))

	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return fmt.Errorf("error removing key %s", key)
	}
	return nil
}

func (p *PersistentStore) Entries() ([]Entry, error) {
	for _, mu := range p.keyMutexMap {
		mu.Lock()
		defer mu.Unlock()
	}

	keys, err := os.ReadDir(p.storeRoot)
	if err != nil {
		return []Entry{}, fmt.Errorf("error reading contents from path %s", p.storeRoot)
	}

	entries := make([]Entry, 0, len(keys))
	for _, key := range keys {
		name := key.Name()
		val, err := p.getUnsafe(name)
		if err != nil {
			return []Entry{}, err
		}
		entries = append(entries, Entry{name, val})
	}

	return entries, nil
}
