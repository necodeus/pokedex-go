package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time // represents when the entry was created
	val       []byte    // that represents the raw data we're caching
}

type Cache struct {
	data     map[string]cacheEntry // a map[string]cacheEntry
	mutex    sync.Mutex            // a mutex to protect the map across goroutines
	interval time.Duration
}

// a NewCache() function that creates a new cache with a configurable interval (time.Duration)
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
	} // create a new cache
	go cache.reapLoop() // when the cache is created
	return cache
}

// a cache.Add() method that adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheEntry{createdAt: time.Now(), val: val} // add the entry to the map
}

// a cache.Get() method that gets an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exists := c.data[key] // get the entry from the map

	// the entry was not found
	if !exists {
		return nil, false
	}

	// the entry was found
	return entry.val, true
}

// a cache.reapLoop() method that is called when the cache is created (by the NewCache function)
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// a loop that runs every interval
	for range ticker.C {
		c.mutex.Lock()
		for key, entry := range c.data {
			// if the entry is older than the interval
			if time.Since(entry.createdAt) > c.interval {
				delete(c.data, key) // delete the entry from the map
			}
		}
		c.mutex.Unlock()
	}
}
