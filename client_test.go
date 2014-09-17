package main

import (
	"log"
	"testing"
	"time"
)

var (
	c   *Client
	err error

	dsn       = "localhost:9876"
	cacheItem = &CacheItem{Key: "some key", Value: "some value"}
)

func init() {
	c, err = NewClient(dsn, time.Millisecond*500)
	if err != nil {
		log.Fatal(err)
	}
}

func TestColdGet(t *testing.T) {
	item, _ := c.Get(cacheItem.Key)
	if item != nil {
		t.Errorf("Cache key should not exist: %s\n", cacheItem.Key)
	}
}

func TestPut(t *testing.T) {
	_, err := c.Put(cacheItem)
	if err != nil {
		t.Error(err)
	}
}

func TestWarmGet(t *testing.T) {
	item, _ := c.Get(cacheItem.Key)
	if item == nil {
		t.Errorf("Cache key should exist: %s\n", cacheItem.Key)
	}
	if item.Value != cacheItem.Value {
		t.Errorf("Cache value expected %s got %s\n", cacheItem.Value, item.Value)
	}
}

func TestDelete(t *testing.T) {
	_, err := c.Delete(cacheItem.Key)
	if err != nil {
		t.Error(err)
	}

	item, _ := c.Get(cacheItem.Key)
	if item != nil {
		t.Errorf("Cache key should not exist: %s\n", cacheItem.Key)
	}
}

func TestClear(t *testing.T) {
	_, err := c.Clear()
	if err != nil {
		t.Error(err)
	}
}

func TestStats(t *testing.T) {
	stats, err := c.Stats()
	if err != nil {
		t.Error(err)
	}
	if stats.Get != 1 {
		t.Errorf("Get: expected 1, got %d\n", stats.Get)
	}
	if stats.Put != 1 {
		t.Errorf("Put: expected 1, got %d\n", stats.Put)
	}
	if stats.Delete != 1 {
		t.Errorf("Delete: expected 1, got %d\n", stats.Delete)
	}
	if stats.Clear != 1 {
		t.Errorf("Clear: expected 1, got %d\n", stats.Clear)
	}
}
