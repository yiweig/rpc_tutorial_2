package main

import (
	"log"
	"testing"
	"time"
	"math"
	"fmt"
)

var (
	c   *Client
	err error

	dsn = "localhost:9876"
	cacheItem = &CacheItem{Key: "some key", Value: "some value"}
)

func init() {
	c, err = NewClient(dsn, time.Millisecond * 500)
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
	maxCount := int(math.Inf(-1))

	if stats.Get * 3 != c.getCount {
		t.Errorf("Get: expected %d, got %d\n", stats.Get, c.getCount)
	}
	if stats.Put != c.putCount {
		t.Errorf("Put: expected %d, got %d\n", stats.Put, c.putCount)
	}
	if stats.Delete != c.deleteCount {
		t.Errorf("Delete: expected %d, got %d\n", stats.Delete, c.deleteCount)
	}

	fmt.Println(stats.Get)
	fmt.Println(stats.Put)
	fmt.Println(stats.Delete)

	fmt.Println(c.getCount)
	fmt.Println(c.putCount)
	fmt.Println(c.deleteCount)

	if stats.Get > maxCount {
		maxCount = stats.Get
	}
	if stats.Put > maxCount {
		maxCount = stats.Put
	}
	if stats.Delete > maxCount {
		maxCount = stats.Delete
	}

	fmt.Println(maxCount)

	if stats.Clear != maxCount {
		t.Errorf("Clear: expected 1, got %d with %d\n", stats.Clear, maxCount)
	}
}

func TestReset(t *testing.T) {
	reset, err := c.Reset()
	if err != nil {
		t.Error(err)
	}

	if reset == true {
		t.Logf("Reset!")
	} else {
		t.Error(err)
	}
}