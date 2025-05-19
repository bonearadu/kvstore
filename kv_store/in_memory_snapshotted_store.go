package kv_store

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type InMemorySnapshottedStore struct {
	store               *InMemoryStore
	snapshotsRoot       string
	snapshotFreqSeconds int64
}

func NewInMemorySnapshottedStore(snapshotsRoot string, snapshotFreqSeconds int64) *InMemorySnapshottedStore {
	snapshottedStore := &InMemorySnapshottedStore{
		store:               NewInMemoryStore(),
		snapshotsRoot:       snapshotsRoot,
		snapshotFreqSeconds: snapshotFreqSeconds,
	}

	snapshottedStore.loadSnapshot()

	go func() {
		for {
			snapshottedStore.saveSnapshot()
			time.Sleep(time.Duration(snapshotFreqSeconds) * time.Second)
		}
	}()

	return snapshottedStore
}

func (i *InMemorySnapshottedStore) saveSnapshot() error {
	i.store.mu.Lock()
	defer i.store.mu.Unlock()

	snapshots, err := os.ReadDir(i.snapshotsRoot)
	if err != nil {
		log.Printf("Error reading snapshots directory: %v", err)
		return fmt.Errorf("snapshotsRoot path is invalid: %s", i.snapshotsRoot)
	}

	newSnapName := strconv.FormatInt(time.Now().UnixMilli(), 10)
	newSnapPath := path.Join(i.snapshotsRoot, newSnapName)

	log.Printf("Saving snapshot with name %s at path %s", newSnapName, newSnapPath)

	newSnap, err := os.OpenFile(newSnapPath, os.O_CREATE|os.O_WRONLY, fileMode)
	if err != nil {
		log.Printf("Error creating snapshot file: %v", err)
		return fmt.Errorf("failed to create snapshot at path %s", newSnapPath)
	}
	defer newSnap.Close()

	err = gob.NewEncoder(newSnap).Encode(i.store.mapStore)
	if err != nil {
		log.Printf("Error encoding snapshot: %v", err)
		defer os.Remove(newSnapPath)
		return fmt.Errorf("failed to encode state for map %v", i.store.mapStore)
	}

	for _, snap := range snapshots {
		os.Remove(path.Join(i.snapshotsRoot, snap.Name()))
	}

	log.Printf("Snapshot saved successfully at path %s", newSnapPath)
	return nil
}

func (i *InMemorySnapshottedStore) loadSnapshot() error {
	i.store.mu.Lock()
	defer i.store.mu.Unlock()

	snapshots, err := os.ReadDir(i.snapshotsRoot)
	if err != nil {
		return fmt.Errorf("snapshotsRoot path is invalid: %s", i.snapshotsRoot)
	}

	lastSnapTimestamp := int64(0)
	for _, snap := range snapshots {
		timestamp, err := strconv.ParseInt(snap.Name(), 0, 64)
		if err != nil {
			continue
		}
		if timestamp > lastSnapTimestamp {
			lastSnapTimestamp = timestamp
		}
	}

	snapPath := path.Join(i.snapshotsRoot, strconv.FormatInt(lastSnapTimestamp, 10))
	snap, err := os.Open(snapPath)
	if err != nil {
		return fmt.Errorf("failed to read snapshot: %s", snapPath)
	}
	defer snap.Close()

	err = gob.NewDecoder(snap).Decode(&i.store.mapStore)
	if err != nil {
		return fmt.Errorf("failed to decode snapshot from path %s", snapPath)
	}

	return nil
}

func (i *InMemorySnapshottedStore) Put(key string, value string) error {
	return i.store.Put(key, value)
}

func (i *InMemorySnapshottedStore) Get(key string) (string, error) {
	return i.store.Get(key)
}

func (i *InMemorySnapshottedStore) Delete(key string) error {
	return i.store.Delete(key)
}

func (i *InMemorySnapshottedStore) Entries() ([]Entry, error) {
	return i.store.Entries()
}
